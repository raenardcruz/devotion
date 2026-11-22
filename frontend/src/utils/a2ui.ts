export interface A2UIStringValue {
  literalString?: string;
  path?: string;
}

export interface A2UIChildrenRef {
  explicitList?: string[];
}

export interface A2UITabItem {
  title?: A2UIStringValue | string;
  icon?: string;
  child?: string;
}

export interface A2UIKeyValueItem {
  label?: A2UIStringValue | string;
  value?: A2UIStringValue | string;
  icon?: string;
}

export interface A2UIComponentProps {
  children?: A2UIChildrenRef | string[];
  child?: string;
  title?: A2UIStringValue | string;
  subtitle?: A2UIStringValue | string;
  icon?: string;
  tone?: string;
  size?: string;
  variant?: string;
  gap?: string;
  columns?: number;
  ordered?: boolean;
  defaultOpen?: boolean;
  axis?: 'horizontal' | 'vertical';
  tabItems?: A2UITabItem[];
  items?: A2UIKeyValueItem[];
  text?: A2UIStringValue | string;
  quote?: A2UIStringValue | string;
  attribution?: A2UIStringValue | string;
  source?: A2UIStringValue | string;
  ref?: A2UIStringValue | string;
  usageHint?: string;
  label?: A2UIStringValue | string;
  url?: string;
  name?: string;
  [key: string]: unknown;
}

export interface A2UIComponentNode {
  id: string;
  component: Record<string, A2UIComponentProps>;
}

export interface A2UIDataEntry {
  key: string;
  valueString?: string;
  valueNumber?: number;
  valueBoolean?: boolean;
  valueMap?: A2UIDataEntry[];
}

export interface A2UISurfaceState {
  surfaceId: string;
  rootId?: string;
  components: Record<string, A2UIComponentNode>;
  componentOrder: string[];
  dataModel: Record<string, unknown>;
}

export interface ParsedA2UIResult {
  isA2UI: boolean;
  isStreamingPayload: boolean;
  activeSurface: A2UISurfaceState | null;
  surfaces: Record<string, A2UISurfaceState>;
}

export function stripA2UICodeFences(raw: string): string {
  const trimmed = raw.trim();
  if (!trimmed.startsWith('```')) return trimmed;
  const lines = trimmed.split('\n');
  if (lines.length >= 2 && (lines[0] ?? '').trim().startsWith('```')) {
    lines.shift();
  }
  if (lines.length >= 1 && (lines[lines.length - 1] ?? '').trim() === '```') {
    lines.pop();
  }
  return lines.join('\n').trim();
}

export function looksLikeA2UIStream(raw: string): boolean {
  if (!raw) return false;
  const cleaned = stripA2UICodeFences(raw).trimStart();
  if (!cleaned) return false;
  if (
    cleaned.startsWith('{"surfaceUpdate"') ||
    cleaned.startsWith('{"dataModelUpdate"') ||
    cleaned.startsWith('{"beginRendering"') ||
    cleaned.startsWith('{"deleteSurface"') ||
    ((cleaned.startsWith('{') || cleaned.startsWith('[')) &&
      (cleaned.includes('"surfaceUpdate"') ||
        cleaned.includes('"beginRendering"') ||
        cleaned.includes('"dataModelUpdate"') ||
        cleaned.includes('"components"')))
  ) {
    return true;
  }
  return false;
}

function isA2UIMessageObject(obj: Record<string, unknown>): boolean {
  return Boolean(
    obj.surfaceUpdate ||
      obj.dataModelUpdate ||
      obj.beginRendering ||
      obj.deleteSurface ||
      Array.isArray(obj.components)
  );
}

/**
 * Extracts balanced top-level JSON objects from a string, supporting both
 * compact single-line JSONL (including progressive streaming) and pretty-printed
 * multi-line JSON objects.
 */
function extractTopLevelJSONObjects(cleaned: string): string[] {
  const objs: string[] = [];
  let depth = 0;
  let inString = false;
  let escaped = false;
  let start = -1;

  for (let i = 0; i < cleaned.length; i++) {
    const ch = cleaned[i];
    if (inString) {
      if (escaped) {
        escaped = false;
        continue;
      }
      if (ch === '\\') {
        escaped = true;
        continue;
      }
      if (ch === '"') {
        inString = false;
      }
      continue;
    }

    if (ch === '"') {
      inString = true;
    } else if (ch === '{') {
      if (depth === 0) {
        start = i;
      }
      depth++;
    } else if (ch === '}') {
      if (depth > 0) {
        depth--;
        if (depth === 0 && start >= 0) {
          objs.push(cleaned.slice(start, i + 1));
          start = -1;
        }
      }
    }
  }
  return objs;
}

function applyDataEntries(target: Record<string, unknown>, entries: A2UIDataEntry[]) {
  for (const entry of entries) {
    if (!entry || !entry.key) continue;
    if (entry.valueString !== undefined) {
      target[entry.key] = entry.valueString;
    } else if (entry.valueNumber !== undefined) {
      target[entry.key] = entry.valueNumber;
    } else if (entry.valueBoolean !== undefined) {
      target[entry.key] = entry.valueBoolean;
    } else if (Array.isArray(entry.valueMap)) {
      const nested: Record<string, unknown> = {};
      applyDataEntries(nested, entry.valueMap);
      target[entry.key] = nested;
    }
  }
}

export function resolveA2UIString(
  val: A2UIStringValue | string | undefined,
  dataModel: Record<string, unknown>
): string {
  if (!val) return '';
  if (typeof val === 'string') return val;
  if (val.literalString !== undefined) return val.literalString;
  if (val.path) {
    const cleanPath = val.path.replace(/^\/+/, '');
    const parts = cleanPath.split('/').filter(Boolean);
    let current: unknown = dataModel;
    for (const part of parts) {
      if (current && typeof current === 'object' && part in (current as Record<string, unknown>)) {
        current = (current as Record<string, unknown>)[part];
      } else {
        current = undefined;
        break;
      }
    }
    if (current !== undefined && current !== null) {
      return String(current);
    }
  }
  return '';
}

export function extractChildIds(props: A2UIComponentProps | undefined): string[] {
  if (!props) return [];
  const list: string[] = [];
  const pushUnique = (id: unknown) => {
    if (typeof id === 'string' && id && !list.includes(id)) {
      list.push(id);
    }
  };

  pushUnique(props.child);

  if (Array.isArray(props.children)) {
    for (const id of props.children) {
      pushUnique(id);
    }
  } else if (props.children && Array.isArray(props.children.explicitList)) {
    for (const id of props.children.explicitList) {
      pushUnique(id);
    }
  }
  return list;
}

export function parseA2UISurface(raw: string): ParsedA2UIResult {
  const isStreamingPayload = looksLikeA2UIStream(raw);
  if (!isStreamingPayload) {
    return {
      isA2UI: false,
      isStreamingPayload: false,
      activeSurface: null,
      surfaces: {},
    };
  }

  const cleaned = stripA2UICodeFences(raw);
  const surfaces: Record<string, A2UISurfaceState> = {};
  let lastSurfaceId: string | null = null;
  let parsedAny = false;

  const ensureSurface = (surfaceId: string): A2UISurfaceState => {
    if (!surfaces[surfaceId]) {
      surfaces[surfaceId] = {
        surfaceId,
        components: {},
        componentOrder: [],
        dataModel: {},
      };
    }
    lastSurfaceId = surfaceId;
    return surfaces[surfaceId];
  };

  const applyMessage = (msg: Record<string, unknown>) => {
    if (msg.surfaceUpdate && typeof msg.surfaceUpdate === 'object') {
      const su = msg.surfaceUpdate as { surfaceId?: string; components?: A2UIComponentNode[] };
      const surfaceId = su.surfaceId || 'response';
      const surface = ensureSurface(surfaceId);
      if (Array.isArray(su.components)) {
        for (const comp of su.components) {
          if (comp && comp.id && comp.component) {
            if (!surface.components[comp.id]) {
              surface.componentOrder.push(comp.id);
            }
            surface.components[comp.id] = comp;
          }
        }
      }
      parsedAny = true;
    } else if (msg.dataModelUpdate && typeof msg.dataModelUpdate === 'object') {
      const dmu = msg.dataModelUpdate as { surfaceId?: string; contents?: A2UIDataEntry[] };
      const surfaceId = dmu.surfaceId || 'response';
      const surface = ensureSurface(surfaceId);
      if (Array.isArray(dmu.contents)) {
        applyDataEntries(surface.dataModel, dmu.contents);
      }
      parsedAny = true;
    } else if (msg.beginRendering && typeof msg.beginRendering === 'object') {
      const br = msg.beginRendering as { surfaceId?: string; root?: string };
      const surfaceId = br.surfaceId || 'response';
      const surface = ensureSurface(surfaceId);
      if (br.root) {
        surface.rootId = br.root;
      }
      parsedAny = true;
    } else if (msg.deleteSurface && typeof msg.deleteSurface === 'object') {
      const ds = msg.deleteSurface as { surfaceId?: string };
      if (ds.surfaceId && surfaces[ds.surfaceId]) {
        delete surfaces[ds.surfaceId];
      }
      parsedAny = true;
    } else if (Array.isArray(msg.components)) {
      const surfaceId = typeof msg.surfaceId === 'string' && msg.surfaceId ? msg.surfaceId : 'response';
      const surface = ensureSurface(surfaceId);
      for (const comp of msg.components as A2UIComponentNode[]) {
        if (comp && comp.id && comp.component) {
          if (!surface.components[comp.id]) {
            surface.componentOrder.push(comp.id);
          }
          surface.components[comp.id] = comp;
        }
      }
      if (typeof msg.root === 'string' && msg.root) {
        surface.rootId = msg.root;
      }
      parsedAny = true;
    }
  };

  // Support JSON array payloads
  let handledArray = false;
  if (cleaned.startsWith('[')) {
    try {
      const arr = JSON.parse(cleaned);
      if (Array.isArray(arr)) {
        for (const item of arr) {
          if (item && typeof item === 'object' && isA2UIMessageObject(item as Record<string, unknown>)) {
            applyMessage(item as Record<string, unknown>);
          }
        }
        handledArray = parsedAny;
      }
    } catch {
      // Fall through to balanced object scanner for partial/streaming payloads
    }
  }

  if (!handledArray) {
    for (const rawObj of extractTopLevelJSONObjects(cleaned)) {
      let msg: Record<string, unknown>;
      try {
        msg = JSON.parse(rawObj) as Record<string, unknown>;
      } catch {
        continue;
      }
      if (isA2UIMessageObject(msg)) {
        applyMessage(msg);
      }
    }
  }

  let activeSurface: A2UISurfaceState | null = null;
  const candidateSurface = lastSurfaceId ? surfaces[lastSurfaceId] : undefined;
  if (candidateSurface) {
    activeSurface = candidateSurface;
    if (!activeSurface.rootId && activeSurface.components['root']) {
      activeSurface.rootId = 'root';
    }
  }

  return {
    isA2UI: parsedAny || isStreamingPayload,
    isStreamingPayload,
    activeSurface,
    surfaces,
  };
}

/**
 * Recovers readable prose from a broken or truncated A2UI JSON payload so
 * raw JSONL is never dumped directly to the user if an LLM gets cut off.
 */
export function extractA2UIFallbackProse(raw: string): string {
  const matches = [...raw.matchAll(/"literalString"\s*:\s*"((?:\\.|[^"\\])*)"/g)];
  if (matches.length === 0) return '';
  const pieces: string[] = [];
  for (const m of matches) {
    const rawLit = m[1];
    if (!rawLit) continue;
    try {
      const decoded = JSON.parse(`"${rawLit}"`) as string;
      if (decoded.trim()) {
        pieces.push(decoded.trim());
      }
    } catch {
      pieces.push(rawLit);
    }
  }
  return pieces.join('\n\n').trim();
}

/**
 * Converts an A2UI JSONL payload into clean structured Markdown (headings,
 * blockquotes, lists, paragraphs) for print/magazine/export consumers.
 * Returns raw unchanged if it is not an A2UI payload.
 */
export function a2uiToMarkdown(raw: string): string {
  const parsed = parseA2UISurface(raw);
  if (!parsed.isA2UI) {
    return raw;
  }
  if (!parsed.activeSurface || Object.keys(parsed.activeSurface.components).length === 0) {
    return extractA2UIFallbackProse(raw) || raw;
  }

  const surface = parsed.activeSurface;
  const rootId = surface.rootId || (surface.components['root'] ? 'root' : surface.componentOrder[0] || '');
  const visited = new Set<string>();
  const blocks: string[] = [];

  const renderInlineListItem = (id: string): string => {
    if (!id || visited.has(id)) return '';
    visited.add(id);
    const node = surface.components[id];
    if (!node?.component) return '';
    const first = Object.entries(node.component)[0];
    if (!first) return '';
    const [type, props] = first;
    if (type === 'Text') {
      return resolveA2UIString(props.text, surface.dataModel).trim();
    }
    if (type === 'Badge') {
      return resolveA2UIString(props.label || props.text, surface.dataModel).trim();
    }
    const subParts: string[] = [];
    for (const cid of extractChildIds(props)) {
      const txt = renderInlineListItem(cid);
      if (txt) subParts.push(txt);
    }
    return subParts.join(' ');
  };

  const visit = (id: string) => {
    if (!id || visited.has(id)) return;
    visited.add(id);
    const node = surface.components[id];
    if (!node?.component) return;

    const first = Object.entries(node.component)[0];
    if (!first) return;
    const [type, props] = first;

    if (type === 'Text') {
      const txt = resolveA2UIString(props.text, surface.dataModel).trim();
      if (!txt) return;
      const hint = props.usageHint || props.variant || 'body';
      if (hint === 'h1') blocks.push(`## ${txt}`);
      else if (hint === 'h2') blocks.push(`### ${txt}`);
      else if (hint === 'h3' || hint === 'h4') blocks.push(`### ${txt}`);
      else if (hint === 'quote') blocks.push(`> ${txt}`);
      else blocks.push(txt);
      return;
    }

    if (type === 'Quote') {
      const txt = resolveA2UIString(props.text || props.quote, surface.dataModel).trim();
      const attr = resolveA2UIString(props.attribution, surface.dataModel).trim();
      const src = resolveA2UIString(props.source || props.ref, surface.dataModel).trim();
      if (txt) {
        const cleanQuote = txt.replace(/^["“]+|["”]+$/g, '');
        let suffix = '';
        if (attr && src) suffix = ` — *${attr}*, ${src}`;
        else if (attr) suffix = ` — *${attr}*`;
        else if (src) suffix = ` — *${src}*`;
        blocks.push(`> "${cleanQuote}"${suffix}`);
      }
      return;
    }

    if (type === 'KeyValue' && Array.isArray(props.items)) {
      const listLines: string[] = [];
      for (const item of props.items) {
        const lbl = resolveA2UIString(item.label, surface.dataModel).trim();
        const val = resolveA2UIString(item.value, surface.dataModel).trim();
        if (lbl && val) listLines.push(`- **${lbl}**: ${val}`);
        else if (val) listLines.push(`- ${val}`);
      }
      if (listLines.length > 0) blocks.push(listLines.join('\n'));
      return;
    }

    if (type === 'List') {
      const childIds = extractChildIds(props);
      const listLines: string[] = [];
      childIds.forEach((cid, idx) => {
        const itemText = renderInlineListItem(cid);
        if (itemText) {
          listLines.push(props.ordered ? `${idx + 1}. ${itemText}` : `- ${itemText}`);
        }
      });
      if (listLines.length > 0) blocks.push(listLines.join('\n'));
      return;
    }

    if (type === 'Tabs' && Array.isArray(props.tabItems)) {
      for (const tab of props.tabItems) {
        const title = resolveA2UIString(tab.title, surface.dataModel).trim();
        if (title) blocks.push(`### ${title}`);
        if (tab.child) visit(tab.child);
      }
      return;
    }

    if (type === 'Card' || type === 'Accordion') {
      const title = resolveA2UIString(props.title, surface.dataModel).trim();
      if (title) blocks.push(`### ${title}`);
      for (const cid of extractChildIds(props)) {
        visit(cid);
      }
      return;
    }

    if (type === 'Badge' || type === 'Icon' || type === 'Divider') {
      // Omit purely decorative badges/icons/dividers from editorial print prose
      return;
    }

    if (type === 'Link') {
      const lbl = resolveA2UIString(props.label, surface.dataModel).trim();
      if (props.url) {
        blocks.push(`[${lbl || props.url}](${props.url})`);
      }
      return;
    }

    // Column, Row, Grid, or unknown container
    for (const cid of extractChildIds(props)) {
      visit(cid);
    }
  };

  if (rootId && surface.components[rootId]) {
    visit(rootId);
  }
  for (const id of surface.componentOrder) {
    if (!visited.has(id)) {
      visit(id);
    }
  }

  return blocks.join('\n\n').trim() || extractA2UIFallbackProse(raw) || raw;
}

export function extractA2UIPlainText(raw: string): string {
  const parsed = parseA2UISurface(raw);
  if (!parsed.isA2UI) {
    return raw;
  }
  if (!parsed.activeSurface || Object.keys(parsed.activeSurface.components).length === 0) {
    return extractA2UIFallbackProse(raw) || raw;
  }

  const surface = parsed.activeSurface;
  const rootId = surface.rootId || (surface.components['root'] ? 'root' : surface.componentOrder[0] || '');
  const visited = new Set<string>();
  const lines: string[] = [];

  const visit = (id: string) => {
    if (!id || visited.has(id)) return;
    visited.add(id);
    const node = surface.components[id];
    if (!node || !node.component) return;

    const entries = Object.entries(node.component);
    const firstEntry = entries[0];
    if (!firstEntry) return;
    const [type, props] = firstEntry;

    if (props.title) {
      const titleText = resolveA2UIString(props.title, surface.dataModel);
      if (titleText) lines.push(titleText);
    }
    if (props.subtitle) {
      const subText = resolveA2UIString(props.subtitle, surface.dataModel);
      if (subText) lines.push(subText);
    }

    if (type === 'Text') {
      const txt = resolveA2UIString(props.text, surface.dataModel);
      if (txt) lines.push(txt);
    } else if (type === 'Quote') {
      const txt = resolveA2UIString(props.text || props.quote, surface.dataModel);
      const attr = resolveA2UIString(props.attribution, surface.dataModel);
      const src = resolveA2UIString(props.source || props.ref, surface.dataModel);
      if (txt) {
        let q = `"${txt.replace(/^["“]+|["”]+$/g, '')}"`;
        if (attr && src) q += ` — ${attr} (${src})`;
        else if (attr) q += ` — ${attr}`;
        else if (src) q += ` — ${src}`;
        lines.push(q);
      }
    } else if (type === 'Badge') {
      const lbl = resolveA2UIString(props.label || props.text, surface.dataModel);
      if (lbl) lines.push(`[${lbl}]`);
    } else if (type === 'KeyValue' && Array.isArray(props.items)) {
      for (const item of props.items) {
        const lbl = resolveA2UIString(item.label, surface.dataModel);
        const val = resolveA2UIString(item.value, surface.dataModel);
        if (lbl && val) lines.push(`${lbl}: ${val}`);
        else if (val) lines.push(val);
      }
    } else if (type === 'Link') {
      const lbl = resolveA2UIString(props.label, surface.dataModel);
      if (lbl || props.url) lines.push(`${lbl || props.url} (${props.url || ''})`);
    } else if (type === 'Tabs' && Array.isArray(props.tabItems)) {
      for (const tab of props.tabItems) {
        const title = resolveA2UIString(tab.title, surface.dataModel);
        if (title) lines.push(title);
        if (tab.child) visit(tab.child);
      }
    }

    for (const childId of extractChildIds(props)) {
      visit(childId);
    }
  };

  if (rootId && surface.components[rootId]) {
    visit(rootId);
  }
  for (const id of surface.componentOrder) {
    if (!visited.has(id)) {
      visit(id);
    }
  }

  return lines.join('\n\n').trim() || extractA2UIFallbackProse(raw) || raw;
}
