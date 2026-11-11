import JSZip from 'jszip';

export interface ChapterData {
  id: string;
  numeral: string;
  title: string;
  content: string[];
}

export interface SectionData {
  id: string;
  title: string;
  chapters: ChapterData[];
}

export interface GenerateEpubOptions {
  title: string;
  author: string;
  era?: string;
  sections: SectionData[];
  coverBlob?: Blob | null;
  onProgress?: (percent: number, statusText: string) => void;
}

/**
 * Escapes characters for strict XML/XHTML.
 */
export function escapeXml(text: string): string {
  if (!text) return '';
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;');
}

/**
 * Formats a markdown paragraph into clean, valid XHTML.
 */
export function formatParagraphXhtml(para: string): { isBlockquote: boolean; html: string } {
  let trimmed = para.trim();
  const isBlockquote = trimmed.startsWith('>');
  if (isBlockquote) {
    trimmed = trimmed.substring(1).trim();
  }

  // XML escape before tag conversions
  let escaped = escapeXml(trimmed);

  // Markdown bold: **text** or __text__
  escaped = escaped.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
  escaped = escaped.replace(/__(.+?)__/g, '<strong>$1</strong>');

  // Markdown italic: *text* or _text_
  escaped = escaped.replace(/\*([^\*]+?)\*/g, '<em>$1</em>');
  escaped = escaped.replace(/_([^_]+?)_/g, '<em>$1</em>');

  // Typographic em-dash
  escaped = escaped.replace(/--|—/g, '&#8212;');

  return { isBlockquote, html: escaped };
}

/**
 * Generates an elegant, high-resolution default book cover on an HTML5 canvas.
 * Returns both a Blob and a data URL for preview.
 */
export async function generateDefaultCover(
  title: string,
  author: string,
  era?: string
): Promise<{ blob: Blob; dataUrl: string }> {
  const width = 1200;
  const height = 1800;

  const canvas = document.createElement('canvas');
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext('2d');

  if (!ctx) {
    throw new Error('Canvas 2D context is not available');
  }

  // Background: warm parchment gradient
  const bgGrad = ctx.createLinearGradient(0, 0, width, height);
  bgGrad.addColorStop(0, '#FAF7F2');
  bgGrad.addColorStop(0.5, '#F7F3EB');
  bgGrad.addColorStop(1, '#EFEAE0');
  ctx.fillStyle = bgGrad;
  ctx.fillRect(0, 0, width, height);

  // Subtle parchment texture / vignette border
  ctx.fillStyle = 'rgba(114, 56, 61, 0.03)';
  ctx.fillRect(40, 40, width - 80, height - 80);

  // Outer border - Burgundy
  ctx.strokeStyle = '#72383D';
  ctx.lineWidth = 4;
  ctx.strokeRect(50, 50, width - 100, height - 100);

  // Inner border - Warm gold / taupe
  ctx.strokeStyle = '#AC9C8D';
  ctx.lineWidth = 1.5;
  ctx.strokeRect(62, 62, width - 124, height - 124);

  // Corner decorations (classical fleurons)
  const drawCorner = (x: number, y: number, angle: number) => {
    ctx.save();
    ctx.translate(x, y);
    ctx.rotate(angle);
    ctx.strokeStyle = '#72383D';
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.moveTo(0, 0);
    ctx.lineTo(24, 0);
    ctx.moveTo(0, 0);
    ctx.lineTo(0, 24);
    ctx.stroke();

    ctx.fillStyle = '#AC9C8D';
    ctx.beginPath();
    ctx.arc(8, 8, 3, 0, Math.PI * 2);
    ctx.fill();
    ctx.restore();
  };

  drawCorner(72, 72, 0);
  drawCorner(width - 72, 72, Math.PI / 2);
  drawCorner(width - 72, height - 72, Math.PI);
  drawCorner(72, height - 72, -Math.PI / 2);

  // Top series label
  ctx.fillStyle = '#72383D';
  ctx.font = 'bold 26px Georgia, "Times New Roman", serif';
  ctx.textAlign = 'center';
  ctx.letterSpacing = '6px';
  ctx.fillText('EARLY CHURCH RESOURCES', width / 2, 170);

  // Divider under top label
  ctx.strokeStyle = '#AC9C8D';
  ctx.lineWidth = 1;
  ctx.beginPath();
  ctx.moveTo(width / 2 - 120, 205);
  ctx.lineTo(width / 2 + 120, 205);
  ctx.stroke();

  // Emblem: Ornate Christian Cross / Chi-Rho
  const emblemY = 480;
  // Circular halo
  ctx.strokeStyle = 'rgba(172, 156, 141, 0.4)';
  ctx.lineWidth = 1.5;
  ctx.beginPath();
  ctx.arc(width / 2, emblemY, 110, 0, Math.PI * 2);
  ctx.stroke();

  ctx.strokeStyle = '#72383D';
  ctx.lineWidth = 2.5;
  ctx.beginPath();
  ctx.arc(width / 2, emblemY, 100, 0, Math.PI * 2);
  ctx.stroke();

  // Cross lines
  ctx.fillStyle = '#72383D';
  // Vertical beam
  ctx.fillRect(width / 2 - 6, emblemY - 70, 12, 140);
  // Horizontal beam
  ctx.fillRect(width / 2 - 50, emblemY - 35, 100, 12);
  // Decorative Cross caps
  ctx.fillRect(width / 2 - 12, emblemY - 76, 24, 6);
  ctx.fillRect(width / 2 - 12, emblemY + 70, 24, 6);
  ctx.fillRect(width / 2 - 56, emblemY - 41, 6, 24);
  ctx.fillRect(width / 2 + 50, emblemY - 41, 6, 24);

  // Small center diamond
  ctx.fillStyle = '#AC9C8D';
  ctx.beginPath();
  ctx.moveTo(width / 2, emblemY - 42);
  ctx.lineTo(width / 2 + 12, emblemY - 30);
  ctx.lineTo(width / 2, emblemY - 18);
  ctx.lineTo(width / 2 - 12, emblemY - 30);
  ctx.closePath();
  ctx.fill();

  // Book Title (Multi-line wrap with dynamic font size)
  let titleFontSize = 62;
  if (title.length > 35) {
    titleFontSize = 46;
  } else if (title.length > 25) {
    titleFontSize = 54;
  }
  ctx.font = `bold ${titleFontSize}px Georgia, "Times New Roman", serif`;
  ctx.fillStyle = '#322D29';
  ctx.letterSpacing = '1px';

  const maxTitleWidth = width - 260;
  const words = title.split(' ');
  const lines: string[] = [];
  let currentLine = words[0] || '';

  for (let i = 1; i < words.length; i++) {
    const testLine = currentLine + ' ' + words[i];
    const metrics = ctx.measureText(testLine);
    if (metrics.width > maxTitleWidth) {
      lines.push(currentLine);
      currentLine = words[i] || '';
    } else {
      currentLine = testLine;
    }
  }
  if (currentLine) {
    lines.push(currentLine);
  }

  const titleStartY = 760;
  const titleLineHeight = titleFontSize * 1.3;
  lines.forEach((line, idx) => {
    ctx.fillText(line, width / 2, titleStartY + idx * titleLineHeight);
  });

  // Author & Era
  const authorY = titleStartY + lines.length * titleLineHeight + 60;

  // Thin dividing line
  ctx.strokeStyle = 'rgba(114, 56, 61, 0.3)';
  ctx.lineWidth = 1;
  ctx.beginPath();
  ctx.moveTo(width / 2 - 80, authorY);
  ctx.lineTo(width / 2 + 80, authorY);
  ctx.stroke();

  // Author Name
  ctx.font = 'italic 40px Georgia, "Times New Roman", serif';
  ctx.fillStyle = '#72383D';
  ctx.letterSpacing = '2px';
  ctx.fillText(author, width / 2, authorY + 60);

  // Era / Dates
  if (era) {
    ctx.font = '24px Georgia, "Times New Roman", serif';
    ctx.fillStyle = '#6E6760';
    ctx.letterSpacing = '3px';
    ctx.fillText(era, width / 2, authorY + 110);
  }

  // Bottom Imprint: DEVOTION
  const bottomY = height - 160;
  ctx.strokeStyle = '#AC9C8D';
  ctx.lineWidth = 1;
  ctx.beginPath();
  ctx.moveTo(width / 2 - 140, bottomY - 35);
  ctx.lineTo(width / 2 + 140, bottomY - 35);
  ctx.stroke();

  ctx.font = 'bold 26px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif';
  ctx.fillStyle = '#72383D';
  ctx.letterSpacing = '8px';
  ctx.fillText('DEVOTION', width / 2, bottomY);

  ctx.font = '16px -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif';
  ctx.fillStyle = '#8C837A';
  ctx.letterSpacing = '4px';
  ctx.fillText('CATHOLIC HERITAGE & THEOLOGY', width / 2, bottomY + 35);

  const dataUrl = canvas.toDataURL('image/jpeg', 0.92);
  const blob = await new Promise<Blob>((resolve, reject) => {
    canvas.toBlob(
      b => {
        if (b) resolve(b);
        else reject(new Error('Failed to generate cover image blob'));
      },
      'image/jpeg',
      0.92
    );
  });

  return { blob, dataUrl };
}

/**
 * Normalizes user-uploaded cover image to standard JPEG for Kindle & e-reader compatibility.
 */
export async function processUploadedCover(
  file: File | Blob
): Promise<{ blob: Blob; dataUrl: string }> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onerror = () => reject(new Error('Failed to read image file'));
    reader.onload = () => {
      const dataUrl = reader.result as string;
      const img = new Image();
      img.onerror = () => reject(new Error('Failed to load image for processing'));
      img.onload = () => {
        // Create canvas to normalize format to JPEG and cap dimensions
        const maxDim = 2400;
        let w = img.width;
        let h = img.height;

        if (w > maxDim || h > maxDim) {
          if (w > h) {
            h = Math.round((h * maxDim) / w);
            w = maxDim;
          } else {
            w = Math.round((w * maxDim) / h);
            h = maxDim;
          }
        }

        const canvas = document.createElement('canvas');
        canvas.width = w;
        canvas.height = h;
        const ctx = canvas.getContext('2d');
        if (!ctx) {
          reject(new Error('Canvas context unavailable'));
          return;
        }

        // Fill white background for transparency in PNGs
        ctx.fillStyle = '#FFFFFF';
        ctx.fillRect(0, 0, w, h);
        ctx.drawImage(img, 0, 0, w, h);

        const outDataUrl = canvas.toDataURL('image/jpeg', 0.92);
        canvas.toBlob(
          blob => {
            if (blob) {
              resolve({ blob, dataUrl: outDataUrl });
            } else {
              reject(new Error('Failed to convert image to JPEG'));
            }
          },
          'image/jpeg',
          0.92
        );
      };
      img.src = dataUrl;
    };
    reader.readAsDataURL(file);
  });
}

/**
 * Stylesheet content for EPUB books.
 */
const EPUB_STYLESHEET = `
@charset "UTF-8";

@page {
  margin: 5% 5%;
}

body {
  margin: 0;
  padding: 0;
  font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif;
  font-size: 1.05em;
  line-height: 1.6;
  color: #24201D;
  text-align: justify;
}

/* Title Page */
body.title-page {
  text-align: center;
  padding: 10% 5%;
}

.title-container {
  margin: 0 auto;
  max-width: 90%;
}

.series-title {
  font-size: 0.9em;
  letter-spacing: 0.25em;
  text-transform: uppercase;
  color: #72383D;
  font-weight: bold;
  margin-bottom: 1.5em;
}

.title-rule {
  width: 60px;
  height: 1px;
  background-color: #AC9C8D;
  margin: 1.5em auto;
}

.book-title {
  font-size: 2em;
  line-height: 1.25;
  color: #322D29;
  margin: 0.5em 0;
  font-weight: normal;
}

.book-author {
  font-size: 1.3em;
  font-style: italic;
  color: #72383D;
  margin: 0.5em 0 0.2em 0;
  font-weight: normal;
}

.book-date {
  font-size: 0.95em;
  color: #786F66;
  margin: 0.2em 0 1.5em 0;
}

.title-cross {
  font-size: 1.5em;
  color: #72383D;
  margin: 1em 0;
}

.publisher {
  font-size: 0.9em;
  font-weight: bold;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: #72383D;
  margin-top: 3em;
  margin-bottom: 0.2em;
}

.publisher-sub {
  font-size: 0.75em;
  letter-spacing: 0.15em;
  color: #8C837A;
  margin: 0;
  text-transform: uppercase;
}

/* Content Headings */
h1.section-title {
  font-size: 1.6em;
  font-weight: bold;
  text-align: center;
  color: #72383D;
  margin-top: 2em;
  margin-bottom: 1.5em;
  padding-bottom: 0.4em;
  border-bottom: 1px solid #D1C7BD;
  page-break-after: avoid;
  break-after: avoid;
}

.chapter-container {
  page-break-before: always;
  break-before: page;
  margin-top: 2em;
  padding-top: 1em;
}

.chapter-header {
  text-align: center;
  margin-bottom: 1.8em;
  page-break-after: avoid;
  break-after: avoid;
}

.chapter-numeral {
  font-size: 0.85em;
  letter-spacing: 0.25em;
  text-transform: uppercase;
  color: #72383D;
  font-weight: bold;
  margin-bottom: 0.4em;
}

.chapter-title {
  font-size: 1.35em;
  line-height: 1.3;
  color: #322D29;
  margin: 0.3em 0 0.6em 0;
  font-weight: normal;
}

.chapter-ornament {
  font-size: 1em;
  color: #AC9C8D;
  margin-top: 0.5em;
}

/* Paragraphs & Typography */
p {
  margin: 0 0 0.8em 0;
  text-indent: 1.5em;
}

p.first-para {
  text-indent: 0;
}

blockquote {
  margin: 1.5em 2em;
  padding: 0.5em 1em;
  border-left: 2px solid #72383D;
  font-style: italic;
  color: #5A2C30;
  text-indent: 0;
}

blockquote p {
  text-indent: 0;
  margin-bottom: 0.4em;
}

/* Table of Contents */
nav#toc ol {
  list-style-type: none;
  padding-left: 0;
}

nav#toc ol li {
  margin-bottom: 0.6em;
}

nav#toc ol li a {
  text-decoration: none;
  color: #72383D;
}

nav#toc ol li ol {
  padding-left: 1.5em;
  margin-top: 0.4em;
}
`;

/**
 * Builds and packages a complete, valid EPUB 3 & EPUB 2 compatible book.
 */
export async function generateEpub(options: GenerateEpubOptions): Promise<Blob> {
  const { title, author, era, sections, onProgress } = options;

  onProgress?.(5, 'Preparing document...');

  const zip = new JSZip();
  const bookId = `urn:uuid:${crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).substring(2)}`;
  const modifiedDate = new Date().toISOString().replace(/\.\d+Z$/, 'Z');

  // 1. mimetype (MUST BE FIRST & UNCOMPRESSED)
  zip.file('mimetype', 'application/epub+zip', { compression: 'STORE' });

  // 2. META-INF/container.xml
  const containerXml = `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`;
  zip.file('META-INF/container.xml', containerXml);

  // 3. Cover Image
  onProgress?.(15, 'Preparing cover artwork...');
  let coverBlob = options.coverBlob;
  if (!coverBlob) {
    const generated = await generateDefaultCover(title, author, era);
    coverBlob = generated.blob;
  }
  const coverBuffer = await coverBlob.arrayBuffer();
  zip.file('OEBPS/images/cover.jpg', coverBuffer);

  // 4. Stylesheet
  zip.file('OEBPS/css/stylesheet.css', EPUB_STYLESHEET);

  // 5. Cover XHTML
  const coverXhtml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="en" lang="en">
<head>
  <title>Cover</title>
  <meta charset="utf-8"/>
  <link rel="stylesheet" type="text/css" href="css/stylesheet.css"/>
  <style type="text/css">
    @page { margin: 0; padding: 0; }
    body.cover-body { margin: 0; padding: 0; text-align: center; background-color: #FAF7F2; }
    div.cover-wrapper { height: 100vh; display: flex; align-items: center; justify-content: center; }
    img.cover-image { max-height: 100%; max-width: 100%; height: auto; width: auto; object-fit: contain; }
  </style>
</head>
<body class="cover-body">
  <div class="cover-wrapper">
    <img class="cover-image" src="images/cover.jpg" alt="Cover" />
  </div>
</body>
</html>`;
  zip.file('OEBPS/cover.xhtml', coverXhtml);

  // 6. Title Page XHTML
  const titleXhtml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
  <title>${escapeXml(title)}</title>
  <meta charset="utf-8"/>
  <link rel="stylesheet" type="text/css" href="css/stylesheet.css"/>
</head>
<body class="title-page">
  <div class="title-container">
    <p class="series-title">Early Church Resources</p>
    <div class="title-rule"></div>
    <h1 class="book-title">${escapeXml(title)}</h1>
    <h2 class="book-author">${escapeXml(author)}</h2>
    ${era ? `<p class="book-date">${escapeXml(era)}</p>` : ''}
    <div class="title-cross">&#10013;</div>
    <div class="title-rule"></div>
    <p class="publisher">Devotion</p>
    <p class="publisher-sub">Catholic Living &amp; Prayer</p>
  </div>
</body>
</html>`;
  zip.file('OEBPS/title.xhtml', titleXhtml);

  // 7. Chapters & Sections Chunking
  onProgress?.(30, 'Formatting chapters and texts...');

  interface ChunkFile {
    filename: string;
    id: string;
    title: string;
    sectionTitle: string;
    chapters: ChapterData[];
  }

  const chunkFiles: ChunkFile[] = [];
  const CHUNK_SIZE = 40; // Max chapters per XHTML file to keep performance snappy

  let fileIndex = 1;
  for (const section of sections) {
    if (section.chapters.length === 0) continue;

    for (let c = 0; c < section.chapters.length; c += CHUNK_SIZE) {
      const slice = section.chapters.slice(c, c + CHUNK_SIZE);
      const isFirstChunk = c === 0;
      const partNum = Math.floor(c / CHUNK_SIZE) + 1;
      const totalParts = Math.ceil(section.chapters.length / CHUNK_SIZE);
      const suffix = totalParts > 1 ? ` (Part ${partNum})` : '';

      chunkFiles.push({
        filename: `chapter_${fileIndex}.xhtml`,
        id: `chunk_${fileIndex}`,
        title: `${section.title}${suffix}`,
        sectionTitle: isFirstChunk ? section.title : '',
        chapters: slice
      });
      fileIndex++;
    }
  }

  // Generate XHTML files for each chunk
  chunkFiles.forEach((chunk) => {
    let contentHtml = '';
    if (chunk.sectionTitle && chunk.sectionTitle !== 'default' && chunk.sectionTitle !== title) {
      contentHtml += `  <h1 class="section-title">${escapeXml(chunk.sectionTitle)}</h1>\n`;
    }

    chunk.chapters.forEach((chap) => {
      contentHtml += `  <section class="chapter-container" id="${escapeXml(chap.id)}">\n`;
      contentHtml += `    <div class="chapter-header">\n`;
      if (chap.numeral) {
        contentHtml += `      <div class="chapter-numeral">Chapter ${escapeXml(chap.numeral)}</div>\n`;
      }
      contentHtml += `      <h2 class="chapter-title">${escapeXml(chap.title)}</h2>\n`;
      contentHtml += `      <div class="chapter-ornament">&#10013;</div>\n`;
      contentHtml += `    </div>\n`;
      contentHtml += `    <div class="chapter-content">\n`;

      let hasFirstPara = false;
      chap.content.forEach(rawPara => {
        if (!rawPara.trim()) return;
        const { isBlockquote, html } = formatParagraphXhtml(rawPara);
        if (isBlockquote) {
          contentHtml += `      <blockquote><p>${html}</p></blockquote>\n`;
        } else {
          const cls = !hasFirstPara ? ' class="first-para"' : '';
          hasFirstPara = true;
          contentHtml += `      <p${cls}>${html}</p>\n`;
        }
      });

      contentHtml += `    </div>\n`;
      contentHtml += `  </section>\n`;
    });

    const fileXhtml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
  <title>${escapeXml(chunk.title)}</title>
  <meta charset="utf-8"/>
  <link rel="stylesheet" type="text/css" href="css/stylesheet.css"/>
</head>
<body>
${contentHtml}
</body>
</html>`;

    zip.file(`OEBPS/${chunk.filename}`, fileXhtml);
  });

  // 8. Navigation Document (EPUB 3: nav.xhtml)
  onProgress?.(55, 'Generating navigation catalog...');

  let navTocHtml = `      <li><a href="title.xhtml">Title Page</a></li>\n`;
  let currentSec = '';
  let insideSubList = false;

  chunkFiles.forEach(chunk => {
    chunk.chapters.forEach(chap => {
      const parentSec = chunk.sectionTitle || 'Content';
      if (parentSec !== currentSec) {
        if (insideSubList) {
          navTocHtml += `        </ol>\n      </li>\n`;
          insideSubList = false;
        }
        if (parentSec !== 'default' && parentSec !== title) {
          navTocHtml += `      <li><span>${escapeXml(parentSec)}</span>\n        <ol>\n`;
          insideSubList = true;
        }
        currentSec = parentSec;
      }

      const chapLabel = chap.numeral ? `Chapter ${chap.numeral}: ${chap.title}` : chap.title;
      if (insideSubList) {
        navTocHtml += `          <li><a href="${chunk.filename}#${escapeXml(chap.id)}">${escapeXml(chapLabel)}</a></li>\n`;
      } else {
        navTocHtml += `      <li><a href="${chunk.filename}#${escapeXml(chap.id)}">${escapeXml(chapLabel)}</a></li>\n`;
      }
    });
  });

  if (insideSubList) {
    navTocHtml += `        </ol>\n      </li>\n`;
  }

  const navXhtml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="en" lang="en">
<head>
  <title>Table of Contents</title>
  <meta charset="utf-8"/>
  <link rel="stylesheet" type="text/css" href="css/stylesheet.css"/>
</head>
<body>
  <nav epub:type="toc" id="toc">
    <h1>Table of Contents</h1>
    <ol>
${navTocHtml}    </ol>
  </nav>
  <nav epub:type="landmarks" hidden="">
    <h2>Landmarks</h2>
    <ol>
      <li><a epub:type="cover" href="cover.xhtml">Cover</a></li>
      <li><a epub:type="toc" href="nav.xhtml">Table of Contents</a></li>
      <li><a epub:type="bodymatter" href="${chunkFiles[0]?.filename || 'title.xhtml'}">Start of Content</a></li>
    </ol>
  </nav>
</body>
</html>`;
  zip.file('OEBPS/nav.xhtml', navXhtml);

  // 9. NCX Document (EPUB 2 compatibility: toc.ncx)
  let ncxPlayOrder = 1;
  let navPointsXml = `    <navPoint id="navPoint-${ncxPlayOrder}" playOrder="${ncxPlayOrder}">
      <navLabel><text>Cover</text></navLabel>
      <content src="cover.xhtml"/>
    </navPoint>\n`;
  ncxPlayOrder++;

  navPointsXml += `    <navPoint id="navPoint-${ncxPlayOrder}" playOrder="${ncxPlayOrder}">
      <navLabel><text>Title Page</text></navLabel>
      <content src="title.xhtml"/>
    </navPoint>\n`;
  ncxPlayOrder++;

  chunkFiles.forEach(chunk => {
    chunk.chapters.forEach(chap => {
      const chapLabel = chap.numeral ? `Chapter ${chap.numeral}: ${chap.title}` : chap.title;
      navPointsXml += `    <navPoint id="navPoint-${ncxPlayOrder}" playOrder="${ncxPlayOrder}">
      <navLabel><text>${escapeXml(chapLabel)}</text></navLabel>
      <content src="${chunk.filename}#${escapeXml(chap.id)}"/>
    </navPoint>\n`;
      ncxPlayOrder++;
    });
  });

  const tocNcx = `<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="${bookId}"/>
    <meta name="dtb:depth" content="2"/>
    <meta name="dtb:totalPageCount" content="0"/>
    <meta name="dtb:maxPageNumber" content="0"/>
  </head>
  <docTitle><text>${escapeXml(title)}</text></docTitle>
  <docAuthor><text>${escapeXml(author)}</text></docAuthor>
  <navMap>
${navPointsXml}  </navMap>
</ncx>`;
  zip.file('OEBPS/toc.ncx', tocNcx);

  // 10. Package Document (OEBPS/content.opf)
  onProgress?.(70, 'Building EPUB package metadata...');

  let manifestXml = `    <item id="cover-image" href="images/cover.jpg" media-type="image/jpeg" properties="cover-image"/>
    <item id="cover" href="cover.xhtml" media-type="application/xhtml+xml"/>
    <item id="nav" href="nav.xhtml" media-type="application/xhtml+xml" properties="nav"/>
    <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"/>
    <item id="css" href="css/stylesheet.css" media-type="text/css"/>
    <item id="title-page" href="title.xhtml" media-type="application/xhtml+xml"/>\n`;

  let spineXml = `    <itemref idref="cover"/>
    <itemref idref="title-page"/>
    <itemref idref="nav"/>\n`;

  chunkFiles.forEach(chunk => {
    manifestXml += `    <item id="${chunk.id}" href="${chunk.filename}" media-type="application/xhtml+xml"/>\n`;
    spineXml += `    <itemref idref="${chunk.id}"/>\n`;
  });

  const contentOpf = `<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="book-id">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:identifier id="book-id">${bookId}</dc:identifier>
    <dc:title>${escapeXml(title)}</dc:title>
    <dc:creator>${escapeXml(author)}</dc:creator>
    <dc:language>en</dc:language>
    <dc:publisher>Devotion</dc:publisher>
    <meta property="dcterms:modified">${modifiedDate}</meta>
    <meta name="cover" content="cover-image"/>
  </metadata>
  <manifest>
${manifestXml}  </manifest>
  <spine toc="ncx">
${spineXml}  </spine>
  <guide>
    <reference type="cover" title="Cover" href="cover.xhtml"/>
    <reference type="toc" title="Table of Contents" href="nav.xhtml"/>
    <reference type="text" title="Beginning" href="${chunkFiles[0]?.filename || 'title.xhtml'}"/>
  </guide>
</package>`;
  zip.file('OEBPS/content.opf', contentOpf);

  // 11. Generate ZIP
  onProgress?.(85, 'Compressing EPUB archive...');

  const epubBlob = await zip.generateAsync(
    {
      type: 'blob',
      mimeType: 'application/epub+zip',
      compression: 'DEFLATE',
      compressionOptions: { level: 6 }
    },
    metadata => {
      const pct = 85 + Math.round((metadata.percent * 14) / 100);
      onProgress?.(pct, `Compressing EPUB archive (${Math.round(metadata.percent)}%)...`);
    }
  );

  onProgress?.(100, 'EPUB generated successfully!');
  return epubBlob;
}

/**
 * Triggers client-side download of a Blob.
 */
export function downloadFile(blob: Blob, filename: string): void {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
