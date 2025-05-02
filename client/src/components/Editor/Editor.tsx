import { EditorState, Transaction } from "prosemirror-state";
import { schema } from "prosemirror-schema-basic";
import {
  ProseMirror,
  ProseMirrorDoc,
  reactKeys,
} from "@handlewithcare/react-prosemirror";
import { useEffect, useState, useCallback } from "react";
import { Schema } from "prosemirror-model";
import { addListNodes } from "prosemirror-schema-list";
import EditorToolbar from "./EditorToolbar";
import { placeholderPlugin } from "./placeholder-plugin";
import "./EditorStyle.css";
import { keymap } from "prosemirror-keymap";
import { baseKeymap } from "prosemirror-commands";
import { useDocumentContext } from "@/context/DocumentContext";
import { history, redo, undo } from "prosemirror-history";
import { collab } from "prosemirror-collab";
import debounce from "lodash.debounce";

// Create an extended schema with list support
const mySchema = new Schema({
  nodes: addListNodes(schema.spec.nodes, "paragraph block*", "block"),
  marks: schema.spec.marks
});

const myPlugins = [
  reactKeys(),
  placeholderPlugin("Write something..."),
  history(),
  collab(),
  keymap({...baseKeymap, "Mod-z": undo, "Mod-Shift-z": redo}),
];

export default function Editor() {
  const { selectedDocument, updateDocument } = useDocumentContext();
  const [editorState, setEditorState] = useState(
    EditorState.create({ 
      schema: mySchema, 
      plugins: myPlugins,
    })
  );

  // Update editor content when selected document changes
  useEffect(() => {
    if (selectedDocument) {
      try {
        // Create a new state with the document content
        const newState = EditorState.create({
          schema: mySchema,
          plugins: myPlugins,
          doc: selectedDocument.content 
            ? mySchema.nodeFromJSON(JSON.parse(selectedDocument.content))
            : undefined
        });
        
        setEditorState(newState);
      } catch (error) {
        console.error("Error parsing document content:", error);
        // Create a default state if content cannot be parsed
        const newState = EditorState.create({
          schema: mySchema,
          plugins: myPlugins,
        });
        
        setEditorState(newState);
      }
    }
  }, [selectedDocument]);

  // Debounced save function to prevent too many API calls
  // eslint-disable-next-line react-hooks/exhaustive-deps
  const debouncedSave = useCallback(
    debounce((docId: string, content: string) => {
      updateDocument(docId, content);
    }, 1000),
    [updateDocument]
  );

  // Handle editor changes and save to backend
  const handleDocChange = useCallback(
    (tr: Transaction) => {
      const newState = editorState.apply(tr);
      setEditorState(newState);
      
      // Save changes if transaction changes the document and we have a selected document
      if (tr.docChanged && selectedDocument) {
        const content = JSON.stringify(newState.doc.toJSON());
        debouncedSave(selectedDocument.id, content);
      }
    },
    [editorState, selectedDocument, debouncedSave]
  );

  // If no document is selected, show a placeholder
  if (!selectedDocument) {
    return (
      <div className="flex-1 p-6 overflow-auto flex justify-center items-center">
        <div className="text-gray-500">
          No document selected or create a new document.
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 p-6 overflow-auto flex flex-col">
      <div className="mb-4 px-4 max-w-3xl mx-auto w-full">
        <h1 className="text-2xl font-bold text-gray-800">{selectedDocument.title}</h1>
        <div className="text-sm text-gray-500">
          Last updated: {new Date(selectedDocument.updated_at).toLocaleString()}
        </div>
      </div>
      <div className="editor-container w-full max-w-3xl mx-auto">
        <div className="editor-wrapper">
          <ProseMirror
            state={editorState}
            dispatchTransaction={handleDocChange}
          >
            <EditorToolbar />
            <ProseMirrorDoc />
          </ProseMirror>
        </div>
      </div>
    </div>
  );
}