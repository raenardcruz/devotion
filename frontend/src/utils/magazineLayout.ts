export type ReadingLengthTier = 'short' | 'medium' | 'long';

export type ScriptureLayoutMode =
  | 'psalm'
  | 'split-lead'
  | 'single-column'
  | 'two-column'
  | 'two-column-long';

export type ExegesisLayoutMode = 'none' | 'brief-card' | 'editorial-essay';

export interface ExegesisSectionBlock {
  id: string;
  heading?: string;
  markdown: string;
  wordCount: number;
  useTwoColumns: boolean;
  isStandfirst?: boolean;
  /** Optional editorial pull quote inserted right before this beat */
  precedingPullQuote?: string;
}

export interface ArticleLayoutInput {
  readingText?: string;
  exegesisRaw?: string;
  isPsalm?: boolean;
  hasImage: boolean;
  includeCommentaries: boolean;
  citationCount: number;
}

export interface ArticleLayoutPlan {
  scriptureWordCount: number;
  exegesisWordCount: number;
  readingLengthTier: ReadingLengthTier;
  scriptureLayoutMode: ScriptureLayoutMode;
  exegesisLayoutMode: ExegesisLayoutMode;
  formattedScriptureHtml: string;
  exegesisMarkdown: string;
  exegesisBlocks: ExegesisSectionBlock[];
  /** Shown between Scripture & Exegesis only when editorially warranted (never redundant) */
  scripturePullQuote: string;
  /** In A4 print, start a long essay cleanly on the next page if current page has <95mm left */
  exegesisStartsNewPageInPrint: boolean;
  /** Label for the exegesis header */
  exegesisKicker: string;
  /** Brief card uses 2-col card grid if sources exist */
  briefCardHasSideSources: boolean;
}

const A4_PRINTABLE_HEIGHT_MM = 275;

export function cleanContextMarkdown(contextStr?: string): string {
  if (!contextStr) return '';
  return contextStr
    .replace(/<details[\s\S]*?<\/details>/gi, '')
    .replace(/\r\n/g, '\n')
    .trim();
}

export function decodeBasicEntities(text: string): string {
  return text
    .replace(/&nbsp;/gi, ' ')
    .replace(/&amp;/gi, '&')
    .replace(/&quot;/gi, '"')
    .replace(/&#39;|&apos;/gi, "'")
    .replace(/&ldquo;/gi, '“')
    .replace(/&rdquo;/gi, '”')
    .replace(/&lsquo;/gi, '‘')
    .replace(/&rsquo;/gi, '’')
    .replace(/&mdash;/gi, '—')
    .replace(/&ndash;/gi, '–')
    .replace(/&hellip;/gi, '…');
}

export function stripHtmlToText(html?: string, stripVerseNumbers = true): string {
  if (!html) return '';
  let cleaned = html;
  if (stripVerseNumbers) {
    cleaned = cleaned.replace(
      /<span[^>]*class="[^"]*\bv\b[^"]*"[^>]*>[\s\S]*?<\/span>/gi,
      ' '
    );
  }
  cleaned = cleaned
    .replace(/<br\s*\/?>/gi, ' ')
    .replace(/<\/(p|div|li|h[1-6]|blockquote)>/gi, ' ')
    .replace(/<[^>]+>/g, ' ');
  return decodeBasicEntities(cleaned).replace(/\s+/g, ' ').trim();
}

export function stripMarkdownToPlainText(md?: string): string {
  if (!md) return '';
  return md
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/[*_`>#]/g, '')
    .replace(/\s+/g, ' ')
    .trim();
}

export function countWords(text: string): number {
  const trimmed = text.trim();
  if (!trimmed) return 0;
  return trimmed.split(/\s+/).filter(Boolean).length;
}

/**
 * Safely enhances Scripture HTML for magazine typography:
 * - Marks a leading API.Bible verse number `<span class="v">N</span>` at the very start
 *   of the first prose paragraph with `mag-v-lead` so the illuminated drop cap wraps
 *   the opening letter of the passage instead of enlarging a verse digit.
 * - Wraps the opening alphabetic character of the first prose paragraph in
 *   `<span class="mag-dropcap-letter">X</span>`.
 */
export function formatScriptureHtml(
  html?: string,
  options?: { isPsalm?: boolean; enableDropCap?: boolean }
): string {
  if (!html) return '';
  const trimmed = html.trim();
  if (!trimmed || options?.isPsalm || !options?.enableDropCap) {
    return trimmed;
  }

  const injectDropCapIntoInner = (inner: string): string | null => {
    const match = inner.match(
      /^(\s*)(<span[^>]*class="([^"]*\bv\b[^"]*)"[^>]*>[^<]*<\/span>\s*)?([“"‘'(]*)([A-Za-zÀ-ÖØ-öø-ÿ])([\s\S]*)$/
    );
    if (!match) return null;
    const [, leadingSpace, verseSpan, verseClass, openQuote, firstLetter, rest] = match;

    let transformedVerse = '';
    if (verseSpan && verseClass) {
      const updatedClass = verseClass.includes('mag-v-lead')
        ? verseClass
        : `${verseClass} mag-v-lead`;
      transformedVerse = verseSpan.replace(
        `class="${verseClass}"`,
        `class="${updatedClass}"`
      );
    }

    const quoteHtml = openQuote
      ? `<span class="mag-dropcap-quote">${openQuote}</span>`
      : '';
    return `${leadingSpace}${transformedVerse}${quoteHtml}<span class="mag-dropcap-letter">${firstLetter}</span>${rest}`;
  };

  // Match paragraph tags and target the first prose paragraph (skipping API.Bible heading classes like s, s1, s2, d, r, ms)
  const headingClassRe = /\b(?:s|s1|s2|s3|ms|ms1|mr|r|d|sp)\b/i;
  const pTagRe = /<p\b([^>]*)>([\s\S]*?)<\/p>/gi;
  let replaced = false;

  const withParagraphs = trimmed.replace(pTagRe, (fullMatch, attrs: string, inner: string) => {
    if (replaced) return fullMatch;
    const classMatch = attrs.match(/class="([^"]*)"/i);
    if (classMatch?.[1] && headingClassRe.test(classMatch[1])) {
      return fullMatch;
    }
    const plainInner = stripHtmlToText(inner, true);
    if (!plainInner) {
      return fullMatch;
    }
    const updated = injectDropCapIntoInner(inner);
    if (!updated) {
      return fullMatch;
    }
    replaced = true;
    return `<p${attrs}>${updated}</p>`;
  });

  if (replaced) {
    return withParagraphs;
  }

  // Fallback for plain text or HTML fragments without <p> tags
  if (!/<p\b/i.test(trimmed)) {
    const updated = injectDropCapIntoInner(trimmed);
    if (updated) return updated;
  }

  return trimmed;
}

const META_OPENER_RE =
  /^(?:in this (?:passage|reading|gospel|excerpt|pericope)|this (?:passage|reading|gospel|excerpt|pericope|text)|today['’]s (?:reading|gospel|passage|liturgy)|these verses|the (?:first|second) reading)\b/i;

function pickBestSentence(
  plainText: string,
  options?: { skipFirst?: boolean; preferMiddle?: boolean }
): string {
  if (!plainText) return '';
  const sentences = (plainText.match(/[^.!?]+[.!?]+/g) || [])
    .map((s) => s.replace(/^[\s“"‘'—–-]+/, '').replace(/[\s”"’']+$/, '').trim())
    .filter(Boolean);

  if (sentences.length === 0) return '';

  const startIdx = options?.skipFirst && sentences.length > 2 ? 1 : 0;
  const candidates = sentences
    .slice(startIdx)
    .map((sentence, offset) => ({ sentence, idx: startIdx + offset }))
    .filter(({ sentence }) => {
      if (sentence.length < 52 || sentence.length > 175) return false;
      // Avoid headings/labels like "Historical Setting: ..."
      const colonIdx = sentence.indexOf(':');
      if (colonIdx > 0 && colonIdx < 32) return false;
      if (META_OPENER_RE.test(sentence)) return false;
      return true;
    });

  if (candidates.length === 0) {
    // Relax meta-opener filter as fallback
    const fallback = sentences
      .slice(startIdx)
      .find((s) => s.length >= 52 && s.length <= 175);
    return fallback || '';
  }

  if (options?.preferMiddle && candidates.length >= 3) {
    return candidates[Math.floor(candidates.length / 2)]?.sentence ?? '';
  }
  return candidates[0]?.sentence ?? '';
}

export function extractCoverDeckQuote(
  readingText?: string,
  exegesisRaw?: string
): string {
  const exegesisPlain = stripMarkdownToPlainText(cleanContextMarkdown(exegesisRaw));
  const fromExegesis = pickBestSentence(exegesisPlain, { preferMiddle: false });
  if (fromExegesis) return fromExegesis;

  const scripturePlain = stripHtmlToText(readingText, true);
  return pickBestSentence(scripturePlain, { preferMiddle: false });
}

interface RawSection {
  heading?: string;
  units: string[];
}

/**
 * Splits a markdown body into indivisible content units (paragraphs, contiguous lists,
 * or contiguous blockquotes) so subdivision never breaks a list or quote in half.
 */
function splitMarkdownIntoUnits(markdown: string): string[] {
  const lines = markdown.split('\n');
  const units: string[] = [];
  let buffer: string[] = [];
  let mode: 'none' | 'list' | 'quote' | 'para' = 'none';

  const flush = () => {
    const joined = buffer.join('\n').trim();
    if (joined) units.push(joined);
    buffer = [];
    mode = 'none';
  };

  for (const rawLine of lines) {
    const trimmed = rawLine.trim();
    if (!trimmed) {
      flush();
      continue;
    }

    const isListLine = /^([\*\-]|\d+\.)\s+/.test(trimmed);
    const isQuoteLine = /^>\s*/.test(trimmed);

    if (isListLine) {
      if (mode !== 'list') flush();
      mode = 'list';
      buffer.push(trimmed);
    } else if (isQuoteLine) {
      if (mode !== 'quote') flush();
      mode = 'quote';
      buffer.push(trimmed);
    } else {
      if (mode !== 'para') flush();
      mode = 'para';
      buffer.push(trimmed);
    }
  }
  flush();
  return units;
}

function cleanHeadingTitle(raw: string): string {
  return raw
    .replace(/^\*\*|\*\*$/g, '')
    .replace(/^__|__$/g, '')
    .replace(/:\s*$/, '')
    .trim();
}

function detectStandaloneHeading(trimmedLine: string): string | null {
  // Explicit markdown heading: #, ##, ###, ####
  const atxMatch = trimmedLine.match(/^#{1,4}\s+(.+)$/);
  if (atxMatch?.[1]) {
    return cleanHeadingTitle(atxMatch[1]);
  }

  // Standalone bold line functioning as a subheading (e.g. **Historical Context** or **1. Literal Sense:**)
  const boldOnlyMatch = trimmedLine.match(/^\*\*([^*\n]{3,84})\*\*$/);
  if (boldOnlyMatch?.[1]) {
    const candidate = boldOnlyMatch[1].trim();
    // Avoid promoting full bold sentences
    if (!/[.!?]$/.test(candidate) && countWords(candidate) <= 10) {
      return cleanHeadingTitle(candidate);
    }
  }

  return null;
}

function countListItems(markdown: string): number {
  return markdown
    .split('\n')
    .filter((l) => /^([\*\-]|\d+\.)\s+/.test(l.trim())).length;
}

/**
 * Subdivides an overly long unheaded/monolithic list of content units into balanced
 * magazine beats (~170-235 words each) so 2-column blocks stay compact on screen
 * and paginate smoothly in A4 print.
 */
function subdivideLongUnits(
  units: string[],
  targetWordsPerBeat = 205
): string[][] {
  const totalWords = countWords(stripMarkdownToPlainText(units.join('\n\n')));
  if (totalWords <= 275 || units.length < 3) {
    return [units];
  }

  const desiredBeats = Math.max(2, Math.round(totalWords / targetWordsPerBeat));
  const balancedTarget = Math.ceil(totalWords / desiredBeats);

  const beats: string[][] = [];
  let current: string[] = [];
  let currentWords = 0;

  for (let i = 0; i < units.length; i++) {
    const u = units[i];
    if (!u) continue;
    const uWords = countWords(stripMarkdownToPlainText(u));
    const remainingUnits = units.length - i;

    if (
      current.length > 0 &&
      currentWords >= balancedTarget * 0.78 &&
      currentWords + uWords > balancedTarget * 1.15 &&
      remainingUnits >= 1
    ) {
      beats.push(current);
      current = [u];
      currentWords = uWords;
    } else {
      current.push(u);
      currentWords += uWords;
    }
  }

  if (current.length > 0) {
    // Merge a tiny trailing orphan beat (< 55 words) into the previous beat if reasonable
    const lastBeat = beats[beats.length - 1];
    if (lastBeat && currentWords < 55) {
      lastBeat.push(...current);
    } else {
      beats.push(current);
    }
  }

  return beats;
}

export function parseExegesisBlocks(
  exegesisMarkdown: string,
  mode: ExegesisLayoutMode
): ExegesisSectionBlock[] {
  if (!exegesisMarkdown || mode === 'none') return [];

  if (mode === 'brief-card') {
    const wc = countWords(stripMarkdownToPlainText(exegesisMarkdown));
    return [
      {
        id: 'brief-0',
        markdown: exegesisMarkdown,
        wordCount: wc,
        useTwoColumns: false,
      },
    ];
  }

  // Parse into sections split by explicit markdown headings or standalone bold headings
  const lines = exegesisMarkdown.split('\n');
  const rawSections: RawSection[] = [];
  let activeHeading: string | undefined;
  let activeBodyLines: string[] = [];

  const pushActiveSection = () => {
    const body = activeBodyLines.join('\n').trim();
    if (body) {
      rawSections.push({
        heading: activeHeading,
        units: splitMarkdownIntoUnits(body),
      });
    } else if (activeHeading) {
      // Heading with no body yet; carry forward only if followed by content
    }
    activeBodyLines = [];
  };

  for (const line of lines) {
    const trimmed = line.trim();
    const detectedHeading = trimmed ? detectStandaloneHeading(trimmed) : null;
    if (detectedHeading) {
      pushActiveSection();
      activeHeading = detectedHeading;
    } else {
      activeBodyLines.push(line);
    }
  }
  pushActiveSection();

  const hasExplicitHeadings = rawSections.some((s) => Boolean(s.heading));
  const blocks: ExegesisSectionBlock[] = [];
  let interludeQuoteUsed = false;

  rawSections.forEach((sec, secIdx) => {
    const sectionMarkdown = sec.units.join('\n\n');
    const sectionWords = countWords(stripMarkdownToPlainText(sectionMarkdown));
    const isIntroBeforeHeadings =
      !sec.heading && secIdx === 0 && hasExplicitHeadings && sectionWords <= 65;

    if (isIntroBeforeHeadings) {
      blocks.push({
        id: `sec-${secIdx}-lead`,
        markdown: sectionMarkdown,
        wordCount: sectionWords,
        useTwoColumns: false,
        isStandfirst: true,
      });
      return;
    }

    const beats = subdivideLongUnits(sec.units);

    beats.forEach((beatUnits, beatIdx) => {
      const beatMd = beatUnits.join('\n\n');
      const beatWords = countWords(stripMarkdownToPlainText(beatMd));
      const beatListItems = countListItems(beatMd);
      const useTwoColumns =
        beatWords >= 68 || (beatWords >= 45 && (beatUnits.length >= 2 || beatListItems >= 3));

      let precedingPullQuote: string | undefined;
      if (beatIdx === 1 && !interludeQuoteUsed && beatWords >= 80) {
        const candidate = pickBestSentence(
          stripMarkdownToPlainText(beatMd),
          { skipFirst: false, preferMiddle: true }
        );
        if (candidate) {
          precedingPullQuote = candidate;
          interludeQuoteUsed = true;
        }
      }

      blocks.push({
        id: `sec-${secIdx}-beat-${beatIdx}`,
        heading: beatIdx === 0 ? sec.heading : undefined,
        markdown: beatMd,
        wordCount: beatWords,
        useTwoColumns,
        precedingPullQuote,
      });
    });
  });

  return blocks;
}

function estimateScriptureSectionHeightMm(
  scriptureWordCount: number,
  mode: ScriptureLayoutMode,
  hasImage: boolean
): number {
  const headerMm = 26;
  if (mode === 'psalm') {
    const imgMm = hasImage ? 72 : 0;
    const textMm = Math.max(36, (scriptureWordCount / 8) * 6.5);
    return headerMm + imgMm + textMm;
  }
  if (mode === 'split-lead') {
    const rightColMm = Math.max(40, (scriptureWordCount / 9) * 6.4);
    const leftImgMm = 68;
    return headerMm + Math.max(rightColMm, leftImgMm) + 8;
  }
  if (mode === 'single-column') {
    const textMm = Math.max(28, (scriptureWordCount / 11) * 6.8);
    return headerMm + textMm + 8;
  }
  if (mode === 'two-column') {
    const imgMm = hasImage ? 82 : 0;
    const textMm = (scriptureWordCount / 18) * 5.9;
    return headerMm + imgMm + textMm + 8;
  }
  // two-column-long
  const imgMm = hasImage ? 72 : 0;
  const textMm = (scriptureWordCount / 20) * 5.5;
  return headerMm + imgMm + textMm + 8;
}

export function computeArticleLayoutPlan(
  input: ArticleLayoutInput
): ArticleLayoutPlan {
  const isPsalm = Boolean(input.isPsalm);
  const scripturePlain = stripHtmlToText(input.readingText, true);
  const scriptureWordCount = countWords(scripturePlain);

  const exegesisMarkdown = input.includeCommentaries
    ? cleanContextMarkdown(input.exegesisRaw)
    : '';
  const exegesisPlain = stripMarkdownToPlainText(exegesisMarkdown);
  const exegesisWordCount = countWords(exegesisPlain);

  const readingLengthTier: ReadingLengthTier =
    scriptureWordCount < 130
      ? 'short'
      : scriptureWordCount <= 320
        ? 'medium'
        : 'long';

  let scriptureLayoutMode: ScriptureLayoutMode;
  if (isPsalm) {
    scriptureLayoutMode = 'psalm';
  } else if (readingLengthTier === 'short') {
    scriptureLayoutMode = input.hasImage ? 'split-lead' : 'single-column';
  } else if (readingLengthTier === 'medium') {
    scriptureLayoutMode = 'two-column';
  } else {
    scriptureLayoutMode = 'two-column-long';
  }

  // Detect whether the exegesis has explicit thematic structure (headings or multi-item lists)
  const exegesisLines = exegesisMarkdown.split('\n');
  const headingCount = exegesisLines.filter((l) =>
    Boolean(detectStandaloneHeading(l.trim()))
  ).length;
  const listCount = countListItems(exegesisMarkdown);
  const hasRichStructure = headingCount >= 1 || listCount >= 4;

  let exegesisLayoutMode: ExegesisLayoutMode;
  if (!exegesisMarkdown || exegesisWordCount === 0) {
    exegesisLayoutMode = 'none';
  } else if (exegesisWordCount <= 155 && !hasRichStructure) {
    exegesisLayoutMode = 'brief-card';
  } else {
    exegesisLayoutMode = 'editorial-essay';
  }

  const exegesisBlocks = parseExegesisBlocks(exegesisMarkdown, exegesisLayoutMode);

  const formattedScriptureHtml = formatScriptureHtml(input.readingText, {
    isPsalm,
    enableDropCap: !isPsalm && scriptureWordCount >= 25,
  });

  // Only show a standalone pull quote beneath Scripture when:
  // - NOT a Psalm (Psalms are already poetic & centered)
  // - Reading is substantial (>= 190 words) OR no exegesis is displayed
  // - Extracted from Scripture itself (so it never duplicates sentence 1 of the exegesis directly below!)
  let scripturePullQuote = '';
  const hasInterludeQuote = exegesisBlocks.some((b) => Boolean(b.precedingPullQuote));
  if (!isPsalm && !hasInterludeQuote && (scriptureWordCount >= 210 || exegesisLayoutMode === 'none')) {
    if (scriptureWordCount >= 120) {
      scripturePullQuote = pickBestSentence(scripturePlain, {
        skipFirst: true,
        preferMiddle: true,
      });
    }
  }

  // Estimate remaining A4 height after Scripture to decide if a long Exegesis Essay
  // should begin cleanly on the next printed page (avoids a cramped 2-line bottom stub).
  let exegesisStartsNewPageInPrint = false;
  if (exegesisLayoutMode === 'editorial-essay' && exegesisWordCount >= 220) {
    const estScriptureMm =
      estimateScriptureSectionHeightMm(
        scriptureWordCount,
        scriptureLayoutMode,
        input.hasImage
      ) + (scripturePullQuote ? 22 : 0);

    const usedOnLastPage = estScriptureMm % A4_PRINTABLE_HEIGHT_MM;
    const remainderMm =
      usedOnLastPage === 0 ? A4_PRINTABLE_HEIGHT_MM : A4_PRINTABLE_HEIGHT_MM - usedOnLastPage;

    if (remainderMm < 95) {
      exegesisStartsNewPageInPrint = true;
    }
  }

  const exegesisKicker =
    exegesisLayoutMode === 'brief-card'
      ? 'Contextual Meditation'
      : exegesisWordCount >= 340
        ? 'In-Depth Exegesis & Commentary'
        : 'Exegesis & Theological Commentary';

  return {
    scriptureWordCount,
    exegesisWordCount,
    readingLengthTier,
    scriptureLayoutMode,
    exegesisLayoutMode,
    formattedScriptureHtml,
    exegesisMarkdown,
    exegesisBlocks,
    scripturePullQuote,
    exegesisStartsNewPageInPrint,
    exegesisKicker,
    briefCardHasSideSources:
      exegesisLayoutMode === 'brief-card' && input.citationCount > 0,
  };
}

export function computePopeQuoteLayout(popeQuoteHtml?: string): {
  isShort: boolean;
  formattedHtml: string;
} {
  const plain = stripHtmlToText(popeQuoteHtml, true);
  const wc = countWords(plain);
  return {
    isShort: wc < 130,
    formattedHtml: formatScriptureHtml(popeQuoteHtml, {
      isPsalm: false,
      enableDropCap: wc >= 25,
    }),
  };
}
