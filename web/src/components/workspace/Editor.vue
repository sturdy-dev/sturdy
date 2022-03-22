<template>
  <div
    class="bg-white border w-full rounded-md shadow-sm focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500 flex flex-row"
  >
    <editor-content :editor="editor" class="flex-1 px-4 py-2" @keydown="keydown" />

    <div v-if="$slots.default" class="flex-none px-4 py-2 flex flex-col justify-end">
      <div>
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Highlight from '@tiptap/extension-highlight'
import Typography from '@tiptap/extension-typography'
import Link from '@tiptap/extension-link'
import { defineComponent } from 'vue'

interface Data {
  editor: Editor | null
}

export default defineComponent({
  components: {
    EditorContent,
  },

  props: {
    modelValue: {
      type: String,
      default: '',
    },
    editable: Boolean,
    placeholder: {
      type: String,
      default: 'Type something...',
    },
  },
  emits: ['updated'],
  data(): Data {
    return {
      editor: null,
      hasFocus: false,
      haveSentAnyUpdates: false,
      lastUpdateIsEmpty: false,
    }
  },

  watch: {
    modelValue(value: string) {
      if (!this.editor) {
        return
      }

      // Ignore if the editor is focused
      if (this.hasFocus) {
        return
      }

      const isSame = this.editor.getHTML() === value

      if (isSame) {
        return
      }

      this.editor.commands.setContent(this.modelValue, false)
    },
  },

  mounted() {
    this.editor = new Editor({
      content: this.modelValue,
      editorProps: {
        attributes: {
          class: 'prose prose-sm max-w-full focus:outline-none',
          style: 'min-height: 60px',
        },
      },

      extensions: [
        StarterKit,
        Placeholder.configure({
          placeholder: this.placeholder,
        }),
        Highlight,
        Typography,
        Link.extend({
          addKeyboardShortcuts() {
            return {
              'Mod-k': () => {
                const url = window.prompt('URL')
                if (!url) {
                  return
                }
                this.editor.chain().focus().setLink({ href: url }).run()
              },
            }
          },
        }),
      ],
      editable: this.editable,
      onUpdate: () => {
        // Send the first update immediately (no debounce)
        if (!this.haveSentAnyUpdates || this.lastUpdateIsEmpty) {
          this.emitUpdated(true, true)
        } else {
          this.emitUpdated(true, false)
        }
        this.haveSentAnyUpdates = true
      },
      onBlur: () => {
        this.hasFocus = false
        this.emitUpdated(true, true)
      },
      onFocus: () => {
        this.hasFocus = true
      },
    })
  },

  beforeUnmount() {
    this.editor?.destroy()
  },
  methods: {
    emitUpdated(isInteractiveUpdate: boolean, shouldSaveImmediately: boolean) {
      if (!this.editor) {
        return
      }

      // Make sure that we don't mistake an empty paragraph for text
      let content = this.editor.getHTML()
      if (content === '<p></p>') {
        content = ''
      }

      this.lastUpdateIsEmpty = content.length === 0

      this.$emit('updated', {
        content: content,
        isInteractiveUpdate: isInteractiveUpdate,
        shouldSaveImmediately: shouldSaveImmediately,
      })
    },
    keydown(e) {
      e.stopPropagation()

      // Don't open the browsers save dialog box
      if (e.key === 's' && e.metaKey) {
        e.preventDefault()
        this.emitUpdated(true, true)
      }
    },
  },
})
</script>

<style scoped>
:deep(.ProseMirror) {
  margin: 0;
  padding: 0;
}

:deep(.ProseMirror > * + *) {
  margin-top: 0.75em;
}

:deep(.ProseMirror p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  float: left;
  color: rgb(156, 163, 175);
  pointer-events: none;
  height: 0;
}

:deep(.ProseMirror ul li > p),
:deep(.ProseMirror ol li > p),
:deep(.ProseMirror ul li),
:deep(.ProseMirror ol li) {
  margin-top: 0;
  margin-bottom: 0;
}

:deep(.ProseMirror a) {
  color: #68cef8;
  cursor: pointer;
}
</style>
