import type { MangaPanelPlan } from '../composables/useStoryToManga';

export interface NormalizedComicTile {
  panel: MangaPanelPlan;
  layoutRow: number;
  colSpan: 2 | 3 | 4 | 6;
  aspectRatio: string;
}

export interface MangaComicLayoutRow {
  rowNumber: number;
  tiles: NormalizedComicTile[];
}

export function aspectRatioForColSpan(colSpan: number): string {
  switch (colSpan) {
    case 6:
      return '16:9';
    case 4:
      return '4:3';
    case 2:
      return '2:3';
    default:
      return '1:1';
  }
}

export function getTileAspectNumeric(aspectRatio?: string, colSpan?: number): number {
  const ratio = aspectRatio || aspectRatioForColSpan(colSpan || 3);
  switch (ratio) {
    case '16:9':
      return 16 / 9;
    case '21:9':
      return 21 / 9;
    case '4:3':
      return 4 / 3;
    case '3:2':
      return 3 / 2;
    case '2:3':
      return 2 / 3;
    case '3:4':
      return 3 / 4;
    case '9:16':
      return 9 / 16;
    default:
      return 1;
  }
}

export function getTileAspectClass(aspectRatio?: string, colSpan?: number): string {
  const ratio = aspectRatio || aspectRatioForColSpan(colSpan || 3);
  switch (ratio) {
    case '16:9':
      return 'aspect-[16/9]';
    case '21:9':
      return 'aspect-[21/9]';
    case '4:3':
      return 'aspect-[4/3]';
    case '3:2':
      return 'aspect-[3/2]';
    case '2:3':
      return 'aspect-[2/3]';
    case '3:4':
      return 'aspect-[3/4]';
    case '9:16':
      return 'aspect-[9/16]';
    default:
      return 'aspect-square';
  }
}

export function getTileColSpanClass(colSpan: number): string {
  switch (colSpan) {
    case 6:
      return 'md:col-span-6';
    case 4:
      return 'md:col-span-4';
    case 2:
      return 'md:col-span-2';
    default:
      return 'md:col-span-3';
  }
}

function defaultComicRowSpans(total: number): Array<Array<2 | 3 | 4 | 6>> {
  switch (total) {
    case 1:
      return [[6]];
    case 2:
      return [[4, 2]];
    case 3:
      return [[6], [3, 3]];
    case 4:
      return [[6], [4, 2], [6]];
    case 5:
      return [[6], [4, 2], [3, 3]];
    case 6:
      return [[6], [4, 2], [2, 2, 2]];
    default: {
      const rows: Array<Array<2 | 3 | 4 | 6>> = [];
      let remaining = total;
      let rowIdx = 0;
      while (remaining > 0) {
        if (remaining >= 3 && rowIdx % 3 === 2) {
          rows.push([2, 2, 2]);
          remaining -= 3;
        } else if (remaining >= 2 && rowIdx % 2 === 1) {
          rows.push([4, 2]);
          remaining -= 2;
        } else if (remaining >= 2) {
          rows.push([3, 3]);
          remaining -= 2;
        } else {
          rows.push([6]);
          remaining--;
        }
        rowIdx++;
      }
      return rows;
    }
  }
}

export function groupPanelsIntoComicRows(panels: MangaPanelPlan[]): MangaComicLayoutRow[] {
  if (!panels || panels.length === 0) return [];
  const sorted = [...panels].sort((a, b) => a.panelNumber - b.panelNumber);

  // Check if panels already carry valid 6-column row metadata
  let valid = true;
  const rowMap = new Map<number, MangaPanelPlan[]>();
  for (const p of sorted) {
    const r = p.layoutRow || 0;
    const c = p.colSpan || 0;
    if (r <= 0 || (c !== 2 && c !== 3 && c !== 4 && c !== 6)) {
      valid = false;
      break;
    }
    const list = rowMap.get(r) || [];
    list.push(p);
    rowMap.set(r, list);
  }

  if (valid) {
    for (const list of rowMap.values()) {
      const sum = list.reduce((acc, item) => acc + (item.colSpan || 0), 0);
      if (sum !== 6) {
        valid = false;
        break;
      }
    }
  }

  if (valid) {
    return Array.from(rowMap.entries())
      .sort((a, b) => a[0] - b[0])
      .map(([rowNumber, list]) => ({
        rowNumber,
        tiles: list.map((p) => {
          const colSpan = (p.colSpan || 3) as 2 | 3 | 4 | 6;
          return {
            panel: p,
            layoutRow: rowNumber,
            colSpan,
            aspectRatio: p.aspectRatio || aspectRatioForColSpan(colSpan),
          };
        }),
      }));
  }

  // Fallback for legacy stories without layout metadata
  const spansList = defaultComicRowSpans(sorted.length);
  const rows: MangaComicLayoutRow[] = [];
  let idx = 0;
  spansList.forEach((spans, rIdx) => {
    const tiles: NormalizedComicTile[] = [];
    for (const span of spans) {
      const p = sorted[idx++];
      if (!p) break;
      tiles.push({
        panel: p,
        layoutRow: rIdx + 1,
        colSpan: span,
        aspectRatio: aspectRatioForColSpan(span),
      });
    }
    if (tiles.length > 0) {
      rows.push({
        rowNumber: rIdx + 1,
        tiles,
      });
    }
  });
  return rows;
}
