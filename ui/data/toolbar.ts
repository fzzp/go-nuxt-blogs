import type { ToolbarNames } from 'md-editor-v3';

export const toolbars: ToolbarNames[] = [
  'bold',
  'underline',
  'italic',
  'strikeThrough',
  '-',
  'title',
  'sub',
  'sup',
  'quote',
  'unorderedList',
  'orderedList',
  'task',
  '-',
  'codeRow',
  'code',
  'link',
  // 'image',
  0, // 自己扩展，第一个，从下标0开始
  'table',
  'mermaid',
  'katex',
  '-',
  'revoke',
  'next',
  'save',
  1, // 第二个自定义工具
  '=',
  'prettier',
  'pageFullscreen',
  'fullscreen',
  'preview',
  'previewOnly',
  // 'htmlPreview',
  'catalog',
  // 'github'
];