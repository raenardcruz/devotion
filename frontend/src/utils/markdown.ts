export interface MarkdownOptions {
  paragraphClass?: string;
  listClass?: string;
  listItemClass?: string;
}

export function renderMarkdown(markdown: string, options?: MarkdownOptions): string {
  if (!markdown) return '';
  
  let html = markdown
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
  
  // Replace headers
  html = html.replace(/^### (.*?)$/gm, '<h5 class="text-xs font-bold text-parchment-primary-dark mt-4 mb-2 font-serif uppercase tracking-wider">$1</h5>');
  html = html.replace(/^## (.*?)$/gm, '<h4 class="text-sm font-bold text-parchment-primary-dark mt-5 mb-2 font-serif">$1</h4>');
  html = html.replace(/^# (.*?)$/gm, '<h3 class="text-base font-bold text-parchment-primary-dark mt-6 mb-3 font-serif">$1</h3>');
  
  // Replace bold
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong class="font-bold text-parchment-primary-dark">$1</strong>');
  html = html.replace(/__(.*?)__/g, '<strong class="font-bold text-parchment-primary-dark">$1</strong>');
  
  // Replace italic
  html = html.replace(/\*(.*?)\*/g, '<em class="italic">$1</em>');
  html = html.replace(/_(.*?)_/g, '<em class="italic">$1</em>');
  
  // Replace code
  html = html.replace(/`(.*?)`/g, '<code class="bg-parchment-neutral-light/70 px-1 py-0.5 rounded font-mono text-xs text-parchment-secondary">$1</code>');
  
  // Replace blockquotes
  html = html.replace(/^&gt; (.*?)$/gm, '<blockquote class="border-l-2 border-parchment-primary/40 pl-3 italic text-parchment-neutral/80 my-2">$1</blockquote>');
  
  // Process line by line for paragraphs and lists
  const lines = html.split('\n');
  let inList = false;
  const processedLines: string[] = [];
  
  const paragraphClass = options?.paragraphClass ?? 'text-sm md:text-[15px] leading-relaxed text-parchment-neutral/80 mb-3';
  const listClass = options?.listClass ?? 'list-disc pl-5 space-y-1.5 my-3';
  const listItemClass = options?.listItemClass ?? 'text-sm md:text-[14px] text-parchment-neutral/85';

  for (let line of lines) {
    const trimmed = line.trim();
    
    // Check for bullet list item
    const listMatch = trimmed.match(/^[\*\-]\s+(.*)$/);
    if (listMatch) {
      if (!inList) {
        processedLines.push(`<ul class="${listClass}">`);
        inList = true;
      }
      processedLines.push(`<li class="${listItemClass}">${listMatch[1]}</li>`);
      continue;
    } else {
      if (inList) {
        processedLines.push('</ul>');
        inList = false;
      }
    }
    
    if (trimmed === '') {
      processedLines.push('');
      continue;
    }
    
    // If it's already an HTML block tag, keep it
    if (trimmed.startsWith('<h') || trimmed.startsWith('<blockquote') || trimmed.startsWith('</blockquote') || trimmed.startsWith('<ul') || trimmed.startsWith('</ul')) {
      processedLines.push(trimmed);
    } else {
      processedLines.push(`<p class="${paragraphClass}">${trimmed}</p>`);
    }
  }
  
  if (inList) {
    processedLines.push('</ul>');
  }
  
  return processedLines.filter(l => l !== '').join('\n');
}
