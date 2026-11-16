import { jsPDF } from 'jspdf';
import type {
  ComicStyle,
  MangaGeneratedPanel,
  MangaPanelPlan,
  MangaStoryAnalysis,
} from '../composables/useStoryToManga';
import {
  groupPanelsIntoComicRows,
  getTileAspectNumeric,
  type MangaComicLayoutRow,
} from './mangaLayout';

export interface MangaA4PdfOptions {
  date: string;
  formattedDate: string;
  style: ComicStyle;
  storyAnalysis: MangaStoryAnalysis | null;
  moralLesson: string;
  panels: MangaGeneratedPanel[];
  getPanelPlan: (panelNumber: number) => MangaPanelPlan | undefined;
  getPanelNarrative: (panel: MangaPanelPlan) => string;
  loadImageElement: (src: string, r2Key?: string) => Promise<HTMLImageElement>;
}

const A4_WIDTH_MM = 210;
const A4_HEIGHT_MM = 297;
const MARGIN_MM = 12;
const CONTENT_WIDTH_MM = A4_WIDTH_MM - MARGIN_MM * 2; // 186mm
const CONTENT_HEIGHT_MM = A4_HEIGHT_MM - MARGIN_MM * 2; // 273mm

// 1240px across 186mm (~169 DPI) keeps comic lettering and artwork crisp
const CANVAS_W = 1240;
const PX_PER_MM = CANVAS_W / CONTENT_WIDTH_MM;

const wrapCanvasText = (
  ctx: CanvasRenderingContext2D,
  text: string,
  maxWidth: number
): string[] => {
  const trimmed = text.trim();
  if (!trimmed) return [];
  const words = trimmed.split(/\s+/);
  const lines: string[] = [];
  let currentLine = '';

  for (const word of words) {
    const testLine = currentLine ? `${currentLine} ${word}` : word;
    if (ctx.measureText(testLine).width <= maxWidth) {
      currentLine = testLine;
    } else {
      if (currentLine) {
        lines.push(currentLine);
      }
      currentLine = word;
    }
  }
  if (currentLine) {
    lines.push(currentLine);
  }
  return lines;
};

const drawRoundedRect = (
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  w: number,
  h: number,
  r: number
) => {
  const radius = Math.min(r, w / 2, h / 2);
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.lineTo(x + w - radius, y);
  ctx.quadraticCurveTo(x + w, y, x + w, y + radius);
  ctx.lineTo(x + w, y + h - radius);
  ctx.quadraticCurveTo(x + w, y + h, x + w - radius, y + h);
  ctx.lineTo(x + radius, y + h);
  ctx.quadraticCurveTo(x, y + h, x, y + h - radius);
  ctx.lineTo(x, y + radius);
  ctx.quadraticCurveTo(x, y, x + radius, y);
  ctx.closePath();
};

const clipRoundedTop = (
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  w: number,
  h: number,
  r: number
) => {
  ctx.beginPath();
  ctx.moveTo(x + r, y);
  ctx.lineTo(x + w - r, y);
  ctx.quadraticCurveTo(x + w, y, x + w, y + r);
  ctx.lineTo(x + w, y + h);
  ctx.lineTo(x, y + h);
  ctx.lineTo(x, y + r);
  ctx.quadraticCurveTo(x, y, x + r, y);
  ctx.closePath();
};

const renderStoryHeaderCanvas = (options: MangaA4PdfOptions): HTMLCanvasElement => {
  const { formattedDate, style, storyAnalysis, moralLesson } = options;

  const measureCanvas = document.createElement('canvas');
  const measureCtx = measureCanvas.getContext('2d')!;

  const cardPadX = 32;
  const cardPadY = 24;
  const innerW = CANVAS_W - cardPadX * 2;

  const eyebrowFont = 'bold 16px "Comic Neue", "Comic Sans MS", sans-serif';
  const titleFont = 'bold 38px "Bangers", "Comic Neue", Impact, sans-serif';
  const settingFont = 'bold 15px "Comic Neue", "Comic Sans MS", sans-serif';
  const lessonLabelFont = 'bold 14px "Comic Neue", "Comic Sans MS", sans-serif';
  const lessonFont = 'italic bold 20px "Comic Neue", "Comic Sans MS", sans-serif';

  const titleText =
    storyAnalysis?.title || `Daily Readings ${style === 'manga' ? 'Manga' : 'Comic'}`;
  measureCtx.font = titleFont;
  const titleLines = wrapCanvasText(measureCtx, titleText, innerW);

  const settingText = storyAnalysis?.setting
    ? `${storyAnalysis.setting.location} · ${storyAnalysis.setting.mood}`
    : '';

  const trimmedLesson = moralLesson.trim();
  measureCtx.font = lessonFont;
  const lessonLines = trimmedLesson
    ? wrapCanvasText(measureCtx, trimmedLesson, innerW - 44)
    : [];

  const eyebrowH = 22;
  const titleLineH = 44;
  const titleBlockH = titleLines.length * titleLineH;
  const settingH = settingText ? 32 : 0;
  const dividerGap = 16;
  const lessonBoxH =
    lessonLines.length > 0 ? 18 + 20 + lessonLines.length * 26 + 14 : 0;

  const cardH =
    cardPadY +
    eyebrowH +
    6 +
    titleBlockH +
    (settingH > 0 ? 8 + settingH : 0) +
    (lessonBoxH > 0 ? dividerGap + lessonBoxH : 0) +
    cardPadY;

  const canvas = document.createElement('canvas');
  canvas.width = CANVAS_W;
  canvas.height = Math.ceil(cardH);
  const ctx = canvas.getContext('2d')!;

  ctx.fillStyle = '#FFFFFF';
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  ctx.fillStyle = '#FFFFFF';
  ctx.strokeStyle = '#1E1B18';
  ctx.lineWidth = 3.5;
  drawRoundedRect(ctx, 2, 2, CANVAS_W - 4, canvas.height - 4, 14);
  ctx.fill();
  ctx.stroke();

  let cursorY = cardPadY;

  ctx.fillStyle = '#72383D';
  ctx.font = eyebrowFont;
  ctx.textAlign = 'left';
  ctx.fillText(
    `CUSTOM PARABLE · INSPIRED BY ${formattedDate.toUpperCase()} READINGS`,
    cardPadX,
    cursorY + 15
  );
  cursorY += eyebrowH + 6;

  ctx.fillStyle = '#111111';
  ctx.font = titleFont;
  titleLines.forEach((line, idx) => {
    ctx.fillText(line, cardPadX, cursorY + 32 + idx * titleLineH);
  });
  cursorY += titleBlockH;

  if (settingText) {
    cursorY += 8;
    ctx.font = settingFont;
    const pillTextW = ctx.measureText(settingText).width;
    const pillW = Math.min(innerW, pillTextW + 26);
    const pillH = 28;
    ctx.fillStyle = '#FAF6EE';
    ctx.strokeStyle = 'rgba(30, 27, 24, 0.4)';
    ctx.lineWidth = 1.5;
    drawRoundedRect(ctx, cardPadX, cursorY, pillW, pillH, 14);
    ctx.fill();
    ctx.stroke();

    ctx.fillStyle = '#322D29';
    ctx.fillText(settingText, cardPadX + 13, cursorY + 19);
    cursorY += settingH;
  }

  if (lessonLines.length > 0) {
    cursorY += 8;
    ctx.strokeStyle = 'rgba(209, 199, 189, 0.7)';
    ctx.lineWidth = 1.5;
    ctx.beginPath();
    ctx.moveTo(cardPadX, cursorY);
    ctx.lineTo(CANVAS_W - cardPadX, cursorY);
    ctx.stroke();
    cursorY += 12;

    const boxX = cardPadX;
    const boxY = cursorY;
    const boxW = innerW;
    const boxH = lessonBoxH;

    ctx.fillStyle = '#FFFDF5';
    ctx.strokeStyle = '#1E1B18';
    ctx.lineWidth = 2.5;
    drawRoundedRect(ctx, boxX, boxY, boxW, boxH, 12);
    ctx.fill();
    ctx.stroke();

    ctx.fillStyle = '#72383D';
    ctx.fillRect(boxX + 1, boxY + 2, 6, boxH - 4);

    ctx.fillStyle = '#72383D';
    ctx.font = lessonLabelFont;
    ctx.fillText("LESSON APPLIED FROM TODAY'S READINGS", boxX + 20, boxY + 22);

    ctx.fillStyle = '#111111';
    ctx.font = lessonFont;
    lessonLines.forEach((line, idx) => {
      ctx.fillText(line, boxX + 20, boxY + 48 + idx * 26);
    });
  }

  return canvas;
};

const renderComicLayoutRowCanvas = async (
  row: MangaComicLayoutRow,
  panelByNumber: Map<number, MangaGeneratedPanel>,
  options: MangaA4PdfOptions
): Promise<HTMLCanvasElement> => {
  const { getPanelNarrative, loadImageElement } = options;
  const tileCount = row.tiles.length;
  const gutter = 20;
  const totalGutterW = Math.max(0, tileCount - 1) * gutter;
  const availGridW = CANVAS_W - totalGutterW;

  const headerBarH = 34;
  const radius = 12;

  const narrativeFont = 'bold 20px "Comic Neue", "Comic Sans MS", sans-serif';
  const dialogueFont = 'italic bold 18px "Comic Neue", "Comic Sans MS", sans-serif';
  const narrativeLineH = 27;
  const dialogueLineH = 24;
  const capPadX = 18;
  const capPadY = 16;

  const measureCanvas = document.createElement('canvas');
  const measureCtx = measureCanvas.getContext('2d')!;

  interface PreparedTile {
    x: number;
    w: number;
    artH: number;
    panelNumber: number;
    aspectRatio: string;
    cameraAngle: string;
    generated?: MangaGeneratedPanel;
    narrativeLines: string[];
    dialogueLines: string[];
    captionH: number;
  }

  const preparedTiles: PreparedTile[] = [];
  let cursorX = 0;
  let maxArtH = 0;
  let maxCaptionH = 0;

  for (let idx = 0; idx < tileCount; idx++) {
    const tile = row.tiles[idx]!;
    const isLast = idx === tileCount - 1;
    const tileW = isLast
      ? CANVAS_W - cursorX
      : Math.round(availGridW * (tile.colSpan / 6));
    const aspectNumeric = getTileAspectNumeric(tile.aspectRatio, tile.colSpan);
    // Cap full-width 16:9 splash height slightly if needed while keeping proportion
    const artH = Math.round(tileW / aspectNumeric);
    if (artH > maxArtH) maxArtH = artH;

    const textInnerW = Math.max(120, tileW - capPadX * 2);
    const narrative = getPanelNarrative(tile.panel);
    const dialogue = (tile.panel.dialogue || '').trim();

    measureCtx.font = narrativeFont;
    const narrativeLines = wrapCanvasText(measureCtx, narrative, textInnerW);

    measureCtx.font = dialogueFont;
    const dialogueLines = dialogue
      ? wrapCanvasText(measureCtx, `“${dialogue}”`, textInnerW - 24)
      : [];

    const narrativeBlockH = narrativeLines.length * narrativeLineH;
    const dialogueBlockH =
      dialogueLines.length > 0 ? dialogueLines.length * dialogueLineH + 22 : 0;
    const captionH =
      capPadY +
      narrativeBlockH +
      (dialogueBlockH > 0 ? 12 + dialogueBlockH : 0) +
      capPadY;

    if (captionH > maxCaptionH) maxCaptionH = captionH;

    preparedTiles.push({
      x: cursorX,
      w: tileW,
      artH,
      panelNumber: tile.panel.panelNumber,
      aspectRatio: tile.aspectRatio,
      cameraAngle: tile.panel.cameraAngle,
      generated: panelByNumber.get(tile.panel.panelNumber),
      narrativeLines,
      dialogueLines,
      captionH,
    });

    cursorX += tileW + gutter;
  }

  const rowH = headerBarH + maxArtH + maxCaptionH + 4;
  const canvas = document.createElement('canvas');
  canvas.width = CANVAS_W;
  canvas.height = Math.ceil(rowH);
  const ctx = canvas.getContext('2d')!;

  ctx.fillStyle = '#FFFFFF';
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  for (const tile of preparedTiles) {
    const tx = tile.x + 2;
    const ty = 2;
    const tw = tile.w - 4;
    const th = rowH - 4;

    ctx.save();
    drawRoundedRect(ctx, tx, ty, tw, th, radius);
    ctx.fillStyle = '#FFFDF5';
    ctx.fill();
    ctx.clip();

    // Header Bar
    ctx.fillStyle = '#1E1B18';
    clipRoundedTop(ctx, tx, ty, tw, headerBarH, radius);
    ctx.fill();

    ctx.fillStyle = '#FFFFFF';
    ctx.font = 'bold 18px "Bangers", "Comic Neue", Impact, sans-serif';
    ctx.textAlign = 'left';
    ctx.fillText(`PANEL ${tile.panelNumber}`, tx + 14, ty + 23);

    // Artwork Area
    const artY = ty + headerBarH;
    ctx.fillStyle = '#F5F0E6';
    ctx.fillRect(tx, artY, tw, maxArtH);

    if (tile.generated?.image) {
      try {
        const img = await loadImageElement(tile.generated.image, tile.generated.r2Key);
        // Cover-fit inside the tile artwork box so there are zero white bars
        const scale = Math.max(tw / img.width, maxArtH / img.height);
        const drawW = img.width * scale;
        const drawH = img.height * scale;
        const drawX = tx + (tw - drawW) / 2;
        const drawY = artY + (maxArtH - drawH) / 2;
        ctx.save();
        ctx.beginPath();
        ctx.rect(tx, artY, tw, maxArtH);
        ctx.clip();
        ctx.drawImage(img, drawX, drawY, drawW, drawH);
        ctx.restore();
      } catch {
        // Skip broken image
      }
    }

    // Divider between artwork and caption strip
    const capY = artY + maxArtH;
    ctx.strokeStyle = '#1E1B18';
    ctx.lineWidth = 2.5;
    ctx.beginPath();
    ctx.moveTo(tx, capY);
    ctx.lineTo(tx + tw, capY);
    ctx.stroke();

    // Narrative Copy
    ctx.fillStyle = '#111111';
    ctx.font = narrativeFont;
    ctx.textAlign = 'left';
    let textCursorY = capY + capPadY + 16;
    for (const line of tile.narrativeLines) {
      ctx.fillText(line, tx + capPadX, textCursorY);
      textCursorY += narrativeLineH;
    }

    // Speech Bubble
    if (tile.dialogueLines.length > 0) {
      const bubbleY = textCursorY + 4;
      const bubbleH = tile.dialogueLines.length * dialogueLineH + 18;
      const bubbleX = tx + capPadX;
      const bubbleW = tw - capPadX * 2;

      ctx.fillStyle = '#FFFFFF';
      ctx.strokeStyle = '#1E1B18';
      ctx.lineWidth = 2;
      drawRoundedRect(ctx, bubbleX, bubbleY, bubbleW, bubbleH, 10);
      ctx.fill();
      ctx.stroke();

      ctx.fillStyle = '#111111';
      ctx.font = dialogueFont;
      tile.dialogueLines.forEach((line, lIdx) => {
        ctx.fillText(line, bubbleX + 12, bubbleY + 23 + lIdx * dialogueLineH);
      });
    }

    ctx.restore();

    // Outer Tile Border
    ctx.strokeStyle = '#1E1B18';
    ctx.lineWidth = 3.5;
    drawRoundedRect(ctx, tx, ty, tw, th, radius);
    ctx.stroke();
  }

  return canvas;
};

export async function exportMangaAsA4Pdf(
  options: MangaA4PdfOptions
): Promise<void> {
  if (options.panels.length === 0) return;

  if (typeof document !== 'undefined' && document.fonts?.ready) {
    await document.fonts.ready;
  }

  const sortedPanels = [...options.panels].sort(
    (a, b) => a.panelNumber - b.panelNumber
  );
  const panelByNumber = new Map<number, MangaGeneratedPanel>();
  const plans: MangaPanelPlan[] = [];

  for (const gp of sortedPanels) {
    panelByNumber.set(gp.panelNumber, gp);
    const plan = options.getPanelPlan(gp.panelNumber);
    if (plan) {
      plans.push(plan);
    } else {
      plans.push({
        panelNumber: gp.panelNumber,
        characters: [],
        sceneDescription: '',
        narrative: '',
        cameraAngle: '',
        visualMood: '',
      });
    }
  }

  const layoutRows = groupPanelsIntoComicRows(plans);

  const pdf = new jsPDF({
    orientation: 'portrait',
    unit: 'mm',
    format: 'a4',
    compress: true,
  });

  const blockGapMm = 4.5;
  const footerReserveMm = 6;
  const maxContentBottomMm = A4_HEIGHT_MM - MARGIN_MM - footerReserveMm;

  let currentYmm = MARGIN_MM;

  // 1. Render & place Story Header Card on Page 1
  const headerCanvas = renderStoryHeaderCanvas(options);
  const headerHeightMm = headerCanvas.height / PX_PER_MM;
  pdf.addImage(
    headerCanvas.toDataURL('image/jpeg', 0.92),
    'JPEG',
    MARGIN_MM,
    currentYmm,
    CONTENT_WIDTH_MM,
    headerHeightMm
  );
  currentYmm += headerHeightMm + blockGapMm;

  // 2. Render & paginate each Comic Layout Row without splitting a row across pages
  for (const row of layoutRows) {
    const rowCanvas = await renderComicLayoutRowCanvas(row, panelByNumber, options);
    const rowHeightMm = rowCanvas.height / PX_PER_MM;

    if (
      currentYmm > MARGIN_MM &&
      currentYmm + rowHeightMm > maxContentBottomMm
    ) {
      pdf.addPage('a4', 'portrait');
      currentYmm = MARGIN_MM;
    }

    const clampedHeightMm = Math.min(rowHeightMm, CONTENT_HEIGHT_MM - footerReserveMm);
    pdf.addImage(
      rowCanvas.toDataURL('image/jpeg', 0.92),
      'JPEG',
      MARGIN_MM,
      currentYmm,
      CONTENT_WIDTH_MM,
      clampedHeightMm
    );
    currentYmm += clampedHeightMm + blockGapMm;
  }

  // 3. Add subtle footer & page numbers to every page
  const totalPages = pdf.getNumberOfPages();
  for (let p = 1; p <= totalPages; p++) {
    pdf.setPage(p);
    pdf.setFont('helvetica', 'normal');
    pdf.setFontSize(8.5);
    pdf.setTextColor(114, 56, 61);
    pdf.text(
      'Catholic Digital Devotional · Story-to-Manga Edition',
      MARGIN_MM,
      A4_HEIGHT_MM - 6
    );
    pdf.setTextColor(92, 81, 73);
    pdf.text(
      `Page ${p} of ${totalPages}`,
      A4_WIDTH_MM - MARGIN_MM,
      A4_HEIGHT_MM - 6,
      { align: 'right' }
    );
  }

  pdf.save(`Devotion_Story_${options.style}_${options.date}.pdf`);
}
