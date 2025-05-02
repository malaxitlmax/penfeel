import { EditorView } from "prosemirror-view";
import { useEditorEventCallback } from "@handlewithcare/react-prosemirror";
import { toggleMark, wrapIn, setBlockType } from "prosemirror-commands";
import { wrapInList } from "prosemirror-schema-list";
import { EditorState, Transaction } from "prosemirror-state";
import "./ToolbarStyle.css";

type CommandCallback = (
  state: EditorState,
  dispatch: ((tr: Transaction) => void) | undefined,
  view: EditorView
) => boolean;

const EditorToolbar = () => {
  const execCommand = useEditorEventCallback((view, callback: CommandCallback) => {
    if (view) {
      callback(view.state, view.dispatch, view);
    }
  });

  const toggleBold = () => {
    execCommand((state: EditorState, dispatch, view: EditorView) => {
      if (!view) return false;
      return toggleMark(state.schema.marks.strong)(state, dispatch, view);
    });
  };

  const toggleItalic = () => {
    execCommand((state: EditorState, dispatch, view: EditorView) => {
      if (!view) return false;
      return toggleMark(state.schema.marks.em)(state, dispatch, view);
    });
  };

  return (
    <div className="editor-toolbar">
      <button className="toolbar-btn" onClick={toggleBold} title="Bold">
        <span className="toolbar-icon">B</span>
      </button>
      <button className="toolbar-btn" onClick={toggleItalic} title="Italic">
        <span className="toolbar-icon"><i>I</i></span>
      </button>
      <div className="toolbar-divider"></div>
    </div>
  );
};

export default EditorToolbar; 