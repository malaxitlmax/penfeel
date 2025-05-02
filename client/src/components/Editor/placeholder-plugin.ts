import { Plugin, PluginKey } from "prosemirror-state";
import { DecorationSet, Decoration } from "prosemirror-view";
import { Node } from "prosemirror-model";

// Creates a plugin that adds placeholder text when content is empty
export function placeholderPlugin(text: string = "Start typing...") {
  return new Plugin({
    key: new PluginKey("placeholder"),
    props: {
      decorations(state) {
        const doc = state.doc;
        
        // If the doc is empty or has only an empty paragraph, show placeholder
        if (isDocEmpty(doc)) {
          const decoration = Decoration.node(0, doc.nodeSize, {
            class: "is-editor-empty",
            "data-placeholder": text,
          });
          
          return DecorationSet.create(doc, [decoration]);
        }
        
        return null;
      },
    },
  });
}

// Helper to check if document is essentially empty
function isDocEmpty(doc: Node): boolean {
  if (doc.childCount === 0) {
    return true;
  }
  
  if (doc.childCount === 1 && doc.firstChild?.isTextblock) {
    const firstChild = doc.firstChild;
    return firstChild.content.size === 0;
  }
  
  return false;
} 