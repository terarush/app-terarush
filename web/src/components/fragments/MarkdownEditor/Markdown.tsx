import MDEditor from '@uiw/react-md-editor';
import '@uiw/react-md-editor/markdown-editor.css';
import '@uiw/react-markdown-preview/markdown.css';
import './MarkdownEditor.css';
import { useTheme } from '@/components/theme-provider';

interface MarkdownEditorProps {
	value: string;
	onChange: (content: string) => void;
}

export function MarkdownEditor({ value, onChange }: MarkdownEditorProps) {
	const { theme } = useTheme();
	// Jika theme adalah 'system', kita perlu tahu apakah sistemnya dark atau light
	const currentTheme = theme === 'system' 
		? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
		: theme;

	return (
		<div className="markdown-editor-wrapper" data-color-mode={currentTheme}>
			<div onKeyDown={(e) => e.stopPropagation()}>
				<MDEditor
					value={value}
					onChange={(val) => onChange(val || '')}
					height={500}
					preview="live"
					visibleDragbar={true}
					hideToolbar={false}
				/>
			</div>
		</div>
	);
}
