import { EditorState } from "prosemirror-state";
import { schema } from "prosemirror-schema-basic";
import {
  ProseMirror,
  ProseMirrorDoc,
  reactKeys,
} from "@handlewithcare/react-prosemirror";
import { useState } from "react";
import { Schema } from "prosemirror-model";
import { addListNodes } from "prosemirror-schema-list";
import EditorToolbar from "./EditorToolbar";
import { placeholderPlugin } from "./placeholder-plugin";
import "./EditorStyle.css";
import { keymap } from "prosemirror-keymap";
import { baseKeymap } from "prosemirror-commands";
// Create an extended schema with list support
const mySchema = new Schema({
  nodes: addListNodes(schema.spec.nodes, "paragraph block*", "block"),
  marks: schema.spec.marks
});

export default function Editor() {
  const [editorState, setEditorState] = useState(
    EditorState.create({ 
      schema: mySchema, 
      plugins: [
        reactKeys(),
        placeholderPlugin("Write something..."),
        keymap(baseKeymap),
      ] 
    })
  );

  return (
    <div className="editor-container">
      <div className="editor-wrapper">
        <ProseMirror
          state={editorState}
          dispatchTransaction={(tr) => {
            setEditorState((s) => s.apply(tr));
          }}
        >
          <EditorToolbar />
          <ProseMirrorDoc />
        </ProseMirror>
      </div>
    </div>
  );
}