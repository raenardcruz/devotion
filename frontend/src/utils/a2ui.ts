export interface A2UIStringValue {
  literalString?: string;
  path?: string;
}

export interface A2UIChildrenRef {
  explicitList?: string[];
}

export interface A2UITabItem {
  title?: A2UIStringValue;
  child?: string;
}

export interface A2UIComponentProps {
  children?: A2UIChildrenRef;
  child?: string;
  title?: A2UIStringValue;
  variant?: string;
  ordered?: boolean;
  axis?: 'horizontal' | 'vertical';
  tabItems?: A2UITabItem[];
  text?: A2UIStringValue;
  usageHint?: string;
  label?: A2UIStringValue;
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
    (cleaned.startsWith('{') &&
      (cleaned.includes('"surfaceUpdate"') ||
        cleaned.includes('"beginRendering"') ||
        cleaned.includes('"dataModelUpdate"')))
  ) {
    return true;
  }
  return false;
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
        dataModel: {},
      };
    }
    lastSurfaceId = surfaceId;
    return surfaces[surfaceId];
  };

  for (const line of cleaned.split('\n')) {
    const trimmed = line.trim();
    if (!trimmed || !trimmed.startsWith('{') || !trimmed.endsWith('}')) {
      continue;
    }

    let msg: Record<string, unknown>;
    try {
      msg = JSON.parse(trimmed) as Record<string, unknown>;
    } catch {
      continue;
    }

    if (msg.surfaceUpdate && typeof msg.surfaceUpdate === 'object') {
      const su = msg.surfaceUpdate as { surfaceId?: string; components?: A2UIComponentNode[] };
      const surfaceId = su.surfaceId || 'response';
      const surface = ensureSurface(surfaceId);
      if (Array.isArray(su.components)) {
        for (const comp of su.components) {
          if (comp && comp.id && comp.component) {
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
    }
  }

  let activeSurface: A2UISurfaceState | null = null;
  const candidateSurface = lastSurfaceId ? surfaces[lastSurfaceId] : undefined;
  if (candidateSurface) {
    activeSurface = candidateSurface;
    // If beginRendering hasn't arrived yet during streaming, default rootId to 'root' if present
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

export function extractA2UIPlainText(raw: string): string {
  const parsed = parseA2UISurface(raw);
  if (!parsed.isA2UI || !parsed.activeSurface) {
    return raw;
  }

  const surface = parsed.activeSurface;
  const rootId = surface.rootId || 'root';
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

    if (type === 'Text') {
      const txt = resolveA2UIString(props.text, surface.dataModel);
      if (txt) lines.push(txt);
    } else if (type === 'Badge') {
      const lbl = resolveA2UIString(props.label, surface.dataModel);
      if (lbl) lines.push(`[${lbl}]`);
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

    if (props.child) {
      visit(props.child);
    }
    if (props.children?.explicitList) {
      for (const childId of props.children.explicitList) {
        visit(childId);
      }
    }
  };

  if (surface.components[rootId]) {
    visit(rootId);
  } else {
    for (const id of Object.keys(surface.components)) {
      visit(id);
    }
  }

  return lines.join('\n\n').trim() || raw;
}
