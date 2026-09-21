// Core Latin Mass (1962 Roman Missal / Divinum Officium) Typesetting & Parsing Engine

export interface MissaDateInfo {
  divinumDate: string; // MM-DD-YYYY
  isoDate: string;     // YYYY-MM-DD
}

export interface MissaItem {
  type: 'major-header' | 'prayer-row';
  id: string;
  title: string;
  latin?: string;
  vernacular?: string;
}

export interface ParsedMissa {
  feast: string;
  commemoration: string;
  rubricsVersion: string;
  items: MissaItem[];
}

export interface BuildHtmlOptions {
  feast: string;
  commemoration?: string;
  rubricsVersion?: string;
  items: MissaItem[];
  dateStr: string;
  autoPrint?: boolean;
}

// Helper to parse and normalize date arguments
export function parseDateArg(arg?: string | null): MissaDateInfo {
  if (!arg || arg === 'today') {
    const now = new Date();
    const mm = String(now.getMonth() + 1).padStart(2, '0');
    const dd = String(now.getDate()).padStart(2, '0');
    const yyyy = now.getFullYear();
    return { divinumDate: `${mm}-${dd}-${yyyy}`, isoDate: `${yyyy}-${mm}-${dd}` };
  }

  // Handle YYYY-MM-DD
  if (/^\d{4}-\d{2}-\d{2}$/.test(arg)) {
    const [yyyy, mm, dd] = arg.split('-');
    return { divinumDate: `${mm}-${dd}-${yyyy}`, isoDate: arg };
  }

  // Handle MM-DD-YYYY or M-D-YYYY
  if (/^\d{1,2}-\d{1,2}-\d{4}$/.test(arg)) {
    const parts = arg.split('-');
    const mm = (parts[0] || '').padStart(2, '0');
    const dd = (parts[1] || '').padStart(2, '0');
    const yyyy = parts[2] || '';
    return { divinumDate: `${mm}-${dd}-${yyyy}`, isoDate: `${yyyy}-${mm}-${dd}` };
  }

  throw new Error(`Invalid date format: "${arg}". Use MM-DD-YYYY or YYYY-MM-DD.`);
}

// Fetch Mass for a given date from Divinum Officium
export async function fetchMissa(divinumDate: string): Promise<string> {
  const params = new URLSearchParams();
  params.append('command', 'praySancta Missa');
  params.append('date', divinumDate);
  params.append('version', 'Rubrics 1960 - 1960');
  params.append('lang2', 'English');

  const controller = typeof AbortController !== 'undefined' ? new AbortController() : null;
  const timeoutId = controller ? setTimeout(() => controller.abort(), 12000) : null;

  try {
    const response = await fetch('https://www.divinumofficium.com/cgi-bin/missa/missa.pl', {
      method: 'POST',
      headers: {
        'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36',
        'Content-Type': 'application/x-www-form-urlencoded'
      },
      body: params.toString(),
      signal: controller ? controller.signal : undefined
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch from Divinum Officium: HTTP ${response.status}`);
    }

    return await response.text();
  } finally {
    if (timeoutId) clearTimeout(timeoutId);
  }
}

export function cleanCell(rawHtml?: string, isLatin: boolean = false): string {
  if (!rawHtml) return '';

  let c = rawHtml
    .replace(/<TD[^>]*>/gi, '')
    .replace(/<\/TD>/gi, '')
    .replace(/<DIV ALIGN=["']?right["']?>[\s\S]*?<\/DIV>/gi, '') // Remove top/next nav links
    .replace(/<A HREF=[^>]*>([\s\S]*?)<\/A>/gi, '$1') // Strip link tags, preserve text
    .replace(/<!--[\s\S]*?-->/g, ''); // Comments

  // Normalize Drop Caps / Initial Capitals
  c = c.replace(/<FONT[^>]*SIZE=["']?\+2["']?[^>]*COLOR=["']?red["']?[^>]*><B><I>([A-Z\u00C0-\u017F])<\/I><\/B><\/FONT>\s*/gi,
    `<span class="prayer-initial">$1</span>`);

  // Subtitle / Prayer names: <FONT SIZE=+1 COLOR="red"><B><I>...</I></B></FONT>
  c = c.replace(/<FONT[^>]*SIZE=["']?\+1["']?[^>]*COLOR=["']?red["']?[^>]*><B><I>([\s\S]*?)<\/I><\/B><\/FONT>/gi,
    `<div class="prayer-subheading">$1</div>`);

  // Sacred Crosses
  c = c.replace(/<span style=[^>]*>✠<\/span>/gi, '✠');
  c = c.replace(/<span style=[^>]*>\+<\/span>/gi, '✠');
  c = c.replace(/✠/g, '<span class="sacred-cross">✠</span>');

  // Red text: Rubrics, Citations, Versicles, Responses, Speakers
  // Rubrics in English only
  c = c.replace(/<FONT[^>]*COLOR=["']?red["']?[^>]*><I>([\s\S]*?)<\/I><\/FONT>/gi, (_match, p1) => {
    const trimmed = p1.trim();
    if (trimmed === 'S.' || trimmed === 'M.' || trimmed === 'C.' || trimmed === 'P.') {
      return `<span class="rubric-speaker">${trimmed}</span> `;
    } else if (trimmed === '℣.' || trimmed === 'V.') {
      return `<span class="rubric-versicle">℣.</span> `;
    } else if (trimmed === '℟.' || trimmed === 'R.') {
      return `<span class="rubric-response">℟.</span> `;
    } else if (/^Ps\b/i.test(trimmed) || /^\d/i.test(trimmed) || /^[A-Z][a-z]+\.?\s*\d/i.test(trimmed)) {
      // Citation (e.g. Ps 36:30-31, Matt 9:9-13, Ioann. 1, 1-14)
      if (isLatin) {
        return ''; // Only in English!
      }
      return `<div class="rubric-citation">${trimmed}</div>`;
    } else {
      // Instructional rubric (e.g. "Iunctis manibus prosequitur:", "Genuflectit dicens:", "Here all kneel.", "of Apostles")
      if (isLatin) {
        return ''; // Only in English!
      }
      return `<div class="rubric-note">${trimmed}</div>`;
    }
  });

  // Strip any remaining red font tags
  c = c.replace(/<FONT[^>]*COLOR=["']?red["']?[^>]*>([\s\S]*?)<\/FONT>/gi, (_match, p1) => {
    if (isLatin) return '';
    return `<div class="rubric-note">${p1}</div>`;
  });

  // Strip other font tags
  c = c.replace(/<\/?FONT[^>]*>/gi, '');

  // Consolidate excess line breaks
  c = c.replace(/(?:<br\s*\/?>\s*){3,}/gi, '<br/><br/>');

  return c.trim();
}

export function parseMissaHtml(html: string): ParsedMissa {
  let feast = 'Sancta Missa';
  let commemoration = '';
  let rubricsVersion = 'Rubrics 1960 - 1960';

  // Extract header information
  const formMatch = html.match(/<FORM[\s\S]*?<TABLE/i);
  if (formMatch) {
    const headerSnippet = formMatch[0];
    const pMatches = headerSnippet.match(/<P ALIGN=["']?CENTER["']?>([\s\S]*?)<\/P>/gi);
    if (pMatches && pMatches.length > 0) {
      const firstP = pMatches[0] || '';
      const scripturaMatch = firstP.match(/Scriptura:[\s\S]*?(?:<\/SPAN>|<\/I>)\s*(?:<FONT[^>]*>)?([^<]+)/i);
      if (scripturaMatch && scripturaMatch[1]) {
        commemoration = scripturaMatch[1].trim();
      }

      const cleanFirst = firstP
        .replace(/<I>[\s\S]*?<\/I>/gi, '')
        .replace(/<[^>]+>/g, ' ')
        .replace(/\s+/g, ' ')
        .trim();
      if (cleanFirst) {
        feast = cleanFirst;
      }
    }

    const vMatch = headerSnippet.match(/Rubrics\s+\d+[\s\w-]*/i);
    if (vMatch) {
      rubricsVersion = vMatch[0].trim();
    }
  }

  // Parse Table rows
  const trMatches = html.match(/<TR[\s\S]*?<\/TR>/gi) || [];
  const items: MissaItem[] = [];

  trMatches.forEach((tr, index) => {
    if (tr.includes('COLSPAN')) {
      const headerText = tr.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim();
      items.push({
        type: 'major-header',
        id: 'hdr-' + index,
        title: headerText
      });
      return;
    }

    const tds = tr.match(/<TD[\s\S]*?<\/TD>/gi) || [];
    if (tds.length === 2 && tds[0] && tds[1]) {
      const td0 = tds[0];
      const td1 = tds[1];
      const lClean = td0.replace(/<TD[^>]*>/gi, '').replace(/<\/TD>/gi, '').replace(/<DIV ALIGN=["']?right["']?>[\s\S]*?<\/DIV>/gi, '').trim();
      const eClean = td1.replace(/<TD[^>]*>/gi, '').replace(/<\/TD>/gi, '').replace(/<DIV ALIGN=["']?right["']?>[\s\S]*?<\/DIV>/gi, '').trim();

      const lRaw = lClean.split(/<br\s*\/?>/i).map(s => s.trim()).filter(Boolean);
      const eRaw = eClean.split(/<br\s*\/?>/i).map(s => s.trim()).filter(Boolean);

      // When line count matches and has multiple lines, align line-by-line!
      if (lRaw.length === eRaw.length && lRaw.length > 1) {
        for (let i = 0; i < lRaw.length; i++) {
          const lineL = lRaw[i] || '';
          const lineE = eRaw[i] || '';
          const cleanL = cleanCell(lineL, true);
          const cleanE = cleanCell(lineE, false);
          if (cleanL || cleanE) {
            const titleMatch = cleanL.match(/<div class="prayer-subheading">([^<]+)<\/div>/i) ||
                               cleanE.match(/<div class="prayer-subheading">([^<]+)<\/div>/i);
            const title = titleMatch && titleMatch[1] ? titleMatch[1].trim() : '';

            items.push({
              type: 'prayer-row',
              id: `p-${index}-${i}`,
              title,
              latin: cleanL,
              vernacular: cleanE
            });
          }
        }
      } else {
        const cleanLatin = cleanCell(td0, true);
        const cleanVernacular = cleanCell(td1, false);
        const titleMatch = cleanLatin.match(/<div class="prayer-subheading">([^<]+)<\/div>/i) ||
                           cleanVernacular.match(/<div class="prayer-subheading">([^<]+)<\/div>/i);
        const title = titleMatch && titleMatch[1] ? titleMatch[1].trim() : '';

        items.push({
          type: 'prayer-row',
          id: 'p-' + index,
          title,
          latin: cleanLatin,
          vernacular: cleanVernacular
        });
      }
    }
  });

  return { feast, commemoration, rubricsVersion, items };
}

export function buildHtmlDocument({ feast, commemoration, rubricsVersion = 'Rubrics 1960 - 1960', items, dateStr, autoPrint = false }: BuildHtmlOptions): string {
  let bodyHtml = '';

  // Filter out trailing empty headers like 'Post Missam' or 'Ante Missam'
  const filteredItems = items.filter((item) => {
    if (item.type === 'major-header') {
      const isAnte = item.title.toLowerCase().includes('ante missam');
      const isPost = item.title.toLowerCase().includes('post missam');
      if (isAnte || isPost) return false;
    }
    return true;
  });

  filteredItems.forEach((item) => {
    if (item.type === 'major-header') {
      bodyHtml += `
      <div class="major-section-banner">
        <div class="byzantine-cross-ornament">☩</div>
        <h2 class="major-section-title">${item.title}</h2>
        <div class="byzantine-knot-line"></div>
      </div>
      `;
      return;
    }

    // Insert traditional Orthodox headpieces before major liturgical divisions
    if (item.title === 'Incipit') {
      bodyHtml += `
      <div class="major-section-banner">
        <div class="byzantine-headpiece">
          <svg width="220" height="20" viewBox="0 0 220 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 10H95M125 10H210" stroke="#C5A059" stroke-width="1.2" stroke-dasharray="4 2"/>
            <circle cx="110" cy="10" r="7" stroke="#8C1D24" stroke-width="1.2" fill="#FFFFFF"/>
            <path d="M110 5.5V14.5M105.5 10H114.5" stroke="#8C1D24" stroke-width="1.2"/>
          </svg>
        </div>
        <div class="byzantine-cross-ornament">☦ ☩ ☦</div>
        <h2 class="major-section-title">MISSA CATECHUMENORUM</h2>
        <div class="major-section-subtitle">THE MASS OF THE CATECHUMENS · PRAYERS AT THE FOOT OF THE ALTAR</div>
        <div class="byzantine-knot-line"></div>
      </div>
      `;
    } else if (item.title === 'Offertorium') {
      bodyHtml += `
      <div class="major-section-banner">
        <div class="byzantine-headpiece">
          <svg width="220" height="20" viewBox="0 0 220 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 10H95M125 10H210" stroke="#C5A059" stroke-width="1.2" stroke-dasharray="4 2"/>
            <circle cx="110" cy="10" r="7" stroke="#8C1D24" stroke-width="1.2" fill="#FFFFFF"/>
            <path d="M110 5.5V14.5M105.5 10H114.5" stroke="#8C1D24" stroke-width="1.2"/>
          </svg>
        </div>
        <div class="byzantine-cross-ornament">☦ ☩ ☦</div>
        <h2 class="major-section-title">MISSA FIDELIUM : OFFERTORIUM</h2>
        <div class="major-section-subtitle">THE MASS OF THE FAITHFUL · OFFERTORY & PREFACE</div>
        <div class="byzantine-knot-line"></div>
      </div>
      `;
    } else if (item.title === 'Canon') {
      bodyHtml += `
      <div class="major-section-banner">
        <div class="byzantine-headpiece">
          <svg width="220" height="20" viewBox="0 0 220 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10 10H95M125 10H210" stroke="#C5A059" stroke-width="1.2" stroke-dasharray="4 2"/>
            <circle cx="110" cy="10" r="7" stroke="#8C1D24" stroke-width="1.2" fill="#FFFFFF"/>
            <path d="M110 5.5V14.5M105.5 10H114.5" stroke="#8C1D24" stroke-width="1.2"/>
          </svg>
        </div>
        <div class="byzantine-cross-ornament">☦ ☩ ☦</div>
        <h2 class="major-section-title">CANON MISSÆ</h2>
        <div class="major-section-subtitle">THE ROMAN CANON · THE SACRED EUCHARISTIC SACRIFICE</div>
        <div class="byzantine-knot-line"></div>
      </div>
      `;
    } else if (item.title === 'Preparatio Communionis' || item.title === 'Agnus Dei') {
      bodyHtml += `
      <div class="major-section-banner">
        <div class="byzantine-cross-ornament">☩</div>
        <h2 class="major-section-title">COMMUNIO</h2>
        <div class="major-section-subtitle">THE COMMUNION OF THE CELEBRANT AND THE FAITHFUL</div>
        <div class="byzantine-knot-line"></div>
      </div>
      `;
    } else if (item.title === 'Postcommunio') {
      bodyHtml += `
      <div class="major-section-banner">
        <div class="byzantine-cross-ornament">☩</div>
        <h2 class="major-section-title">RITUS CONCLUSIONIS</h2>
        <div class="major-section-subtitle">DISMISSAL, BLESSING, AND LAST GOSPEL</div>
        <div class="byzantine-knot-line"></div>
      </div>
      `;
    } else if (item.title === 'Orationes Leonis XIII') {
      bodyHtml += `
      <div class="major-section-banner">
        <div class="byzantine-cross-ornament">☦ ☩ ☦</div>
        <h2 class="major-section-title">PRECES LEONINÆ</h2>
        <div class="major-section-subtitle">PRAYERS AFTER LOW MASS ORDERED BY POPE LEO XIII</div>
        <div class="byzantine-knot-line"></div>
      </div>
      `;
    }

    bodyHtml += `
    <div class="prayer-block">
      <div class="prayer-col col-latin">
        ${item.latin || ''}
      </div>
      <div class="col-divider"></div>
      <div class="prayer-col col-english">
        ${item.vernacular || ''}
      </div>
    </div>
    `;
  });

  // Add final Colophon / Thanksgiving block
  bodyHtml += `
  <div class="colophon-block">
    <div class="byzantine-cross-ornament">✠</div>
    <div class="colophon-title">GRATIARUM ACTIO POST MISSAM</div>
    <div class="colophon-text">Tibi laus, tibi glória, tibi gratiárum áctio in sǽcula sempitérna, O beáta Trínitas. Amen.</div>
    <div class="colophon-motto">☩ AD MAIOREM DEI GLORIAM ☩</div>
  </div>
  `;

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Sancta Missa - ${feast}</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Cinzel:wght@500;700;900&family=Cinzel+Decorative:wght@700;900&family=EB+Garamond:ital,wght@0,400;0,500;0,600;0,700;1,400;1,600&display=swap" rel="stylesheet">
  <style>
    @page {
      background: #FFFFFF;
      size: letter;
      margin: 0.65in 0.6in 0.65in 0.6in;
      @bottom-center {
        content: "— ☩ " counter(page) " ☩ —";
        font-family: 'Cinzel', serif;
        font-size: 8pt;
        color: #8C1D24;
      }
    }

    * {
      box-sizing: border-box;
      -webkit-print-color-adjust: exact;
      print-color-adjust: exact;
    }

    body {
      font-family: 'EB Garamond', Garamond, 'Georgia', serif;
      background-color: #FFFFFF;
      color: #1A1715;
      margin: 0;
      padding: 0;
      font-size: 9.75pt;
      line-height: 1.44;
      text-rendering: optimizeLegibility;
    }

    /* Traditional Frontispiece / Cover Page */
    .title-page {
      min-height: 9.6in;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      align-items: center;
      text-align: center;
      padding: 2.2rem 1.5rem;
      page-break-after: always;
      position: relative;
    }

    .book-border {
      position: absolute;
      inset: 8px;
      border: 3px double #C5A059;
      pointer-events: none;
    }
    .book-border-inner {
      position: absolute;
      inset: 14px;
      border: 1px solid #8C1D24;
      pointer-events: none;
    }

    .corner-cross {
      position: absolute;
      font-size: 13pt;
      color: #8C1D24;
      line-height: 1;
    }
    .corner-tl { top: 3px; left: 3px; }
    .corner-tr { top: 3px; right: 3px; }
    .corner-bl { bottom: 3px; left: 3px; }
    .corner-br { bottom: 3px; right: 3px; }

    .title-head-ornament {
      color: #8C1D24;
      font-size: 24pt;
      margin-bottom: 0.6rem;
      letter-spacing: 0.3em;
    }

    .title-main {
      font-family: 'Cinzel Decorative', 'Cinzel', serif;
      font-size: 26pt;
      font-weight: 900;
      color: #8C1D24;
      letter-spacing: 0.12em;
      margin: 0.4rem 0 0.2rem 0;
      text-transform: uppercase;
    }

    .title-sub {
      font-family: 'Cinzel', serif;
      font-size: 11.5pt;
      font-weight: 700;
      color: #C5A059;
      letter-spacing: 0.25em;
      text-transform: uppercase;
      margin-bottom: 1.25rem;
    }

    .title-divider {
      width: 180px;
      height: 2px;
      background: linear-gradient(90deg, transparent, #8C1D24, #C5A059, #8C1D24, transparent);
      margin: 1.25rem auto;
    }

    .feast-container {
      background: transparent;
      border: 1.5px solid #C5A059;
      box-shadow: 0 0 0 3px rgba(197, 160, 89, 0.15);
      border-radius: 6px;
      padding: 1.4rem 1.8rem;
      margin: 1.5rem auto;
      max-width: 540px;
    }

    .feast-label {
      font-family: 'Cinzel', serif;
      font-size: 8.5pt;
      font-weight: 700;
      letter-spacing: 0.2em;
      color: #8C1D24;
      text-transform: uppercase;
      margin-bottom: 0.45rem;
    }

    .feast-title {
      font-family: 'EB Garamond', serif;
      font-size: 15.5pt;
      font-weight: 700;
      color: #1A1715;
      line-height: 1.35;
      margin-bottom: 0.65rem;
    }

    .feast-commemoration {
      font-family: 'EB Garamond', serif;
      font-size: 10.5pt;
      font-style: italic;
      color: #3C3632;
    }

    .title-footer {
      font-family: 'Cinzel', serif;
      font-size: 8pt;
      letter-spacing: 0.2em;
      color: #8C1D24;
      text-transform: uppercase;
    }

    /* Content Area */
    .content-container {
      position: relative;
    }

    /* Liturgical Section Banners (Orthodox Headpiece Style) */
    .major-section-banner {
      text-align: center;
      margin: 1.3rem 0 0.65rem 0;
      padding: 0.35rem 0;
      page-break-inside: avoid;
      break-inside: avoid;
      page-break-after: avoid;
      break-after: avoid;
    }

    .byzantine-headpiece {
      margin-bottom: 0.15rem;
    }

    .byzantine-cross-ornament {
      color: #8C1D24;
      font-size: 12pt;
      letter-spacing: 0.4em;
      margin-bottom: 0.15rem;
    }

    .major-section-title {
      font-family: 'Cinzel', serif;
      font-size: 11pt;
      font-weight: 700;
      letter-spacing: 0.2em;
      color: #8C1D24;
      margin: 0.15rem 0;
      text-transform: uppercase;
    }

    .major-section-subtitle {
      font-family: 'Cinzel', serif;
      font-size: 7.25pt;
      font-weight: 600;
      letter-spacing: 0.16em;
      color: #C5A059;
      text-transform: uppercase;
      margin-bottom: 0.3rem;
    }

    .byzantine-knot-line {
      width: 100%;
      max-width: 340px;
      height: 1.5px;
      background: linear-gradient(90deg, transparent, #C5A059, #8C1D24, #C5A059, transparent);
      margin: 0.3rem auto;
    }

    /* Column Headers */
    .table-column-headers {
      display: flex;
      border-bottom: 1.5px solid #8C1D24;
      border-top: 1.5px solid #8C1D24;
      padding: 3.5px 0;
      margin-bottom: 0.65rem;
      background: transparent;
      font-family: 'Cinzel', serif;
      font-size: 8pt;
      font-weight: 700;
      letter-spacing: 0.2em;
      color: #8C1D24;
      text-align: center;
      page-break-after: avoid;
      break-after: avoid;
    }
    .table-col-hdr {
      flex: 1;
    }

    /* Symmetrical Parallel Columns */
    .prayer-block {
      display: flex;
      flex-direction: row;
      align-items: stretch;
      page-break-inside: avoid;
      break-inside: avoid;
      margin-bottom: 0.45rem;
      padding-bottom: 0.28rem;
      border-bottom: 0.5px dotted rgba(197, 160, 89, 0.35);
    }

    .prayer-col {
      flex: 1;
      min-width: 0;
      padding: 0 0.5rem;
      font-size: 9.25pt;
      line-height: 1.42;
    }

    .col-latin {
      padding-right: 0.75rem;
    }

    .col-english {
      padding-left: 0.75rem;
    }

    .col-divider {
      width: 1px;
      background: linear-gradient(180deg, rgba(197,160,89,0.25) 0%, rgba(140,29,36,0.3) 50%, rgba(197,160,89,0.25) 100%);
      flex-shrink: 0;
    }

    /* Prayer Subheadings */
    .prayer-subheading {
      font-family: 'Cinzel', serif;
      font-size: 10pt;
      font-weight: 700;
      letter-spacing: 0.12em;
      color: #8C1D24;
      text-transform: uppercase;
      margin: 0.2rem 0 0.2rem 0;
      border-bottom: 1px solid rgba(140, 29, 36, 0.2);
      padding-bottom: 2px;
    }

    /* Illuminated Orthodox-style Initial Capital */
    .prayer-initial {
      font-family: 'Cinzel Decorative', 'Cinzel', serif;
      font-size: 1.45em;
      line-height: 0.9;
      color: #8C1D24;
      font-weight: 700;
      margin-right: 0.5px;
      display: inline-block;
      vertical-align: -0.04em;
    }

    /* Sacred Symbols & Rubrics */
    .sacred-cross {
      color: #8C1D24;
      font-size: 1.15em;
      font-weight: bold;
      vertical-align: -0.05em;
      margin: 0 0.1em;
    }

    .rubric-speaker {
      font-family: 'Cinzel', serif;
      font-weight: 700;
      color: #8C1D24;
      font-size: 8.5pt;
      letter-spacing: 0.05em;
      margin-right: 0.15em;
    }

    .rubric-versicle, .rubric-response {
      color: #8C1D24;
      font-weight: 700;
      font-size: 9.5pt;
      margin-right: 0.15em;
    }

    .rubric-note {
      color: #8C1D24;
      font-style: italic;
      font-size: 8.5pt;
      line-height: 1.35;
      margin: 0.2rem 0;
      padding-left: 0.25rem;
      border-left: 2px solid rgba(140, 29, 36, 0.3);
    }

    .rubric-citation {
      color: #8C1D24;
      font-style: italic;
      font-size: 8.25pt;
      margin: 0.15rem 0;
      font-weight: 600;
      letter-spacing: 0.03em;
    }

    @media print {
      .no-print {
        display: none !important;
      }
    }

    .print-floating-bar {
      position: fixed;
      top: 18px;
      right: 18px;
      z-index: 10000;
      background: rgba(255, 255, 255, 0.95);
      backdrop-filter: blur(10px);
      border: 1px solid #C5A059;
      border-radius: 30px;
      box-shadow: 0 4px 20px rgba(0,0,0,0.15);
      padding: 6px 14px;
      display: flex;
      align-items: center;
      gap: 10px;
    }

    .print-btn {
      background: #8C1D24;
      color: #fff;
      border: none;
      padding: 6px 16px;
      border-radius: 20px;
      font-family: "Cinzel", serif;
      font-size: 8.5pt;
      font-weight: 700;
      letter-spacing: 0.1em;
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 6px;
      transition: background 0.2s;
    }
    .print-btn:hover {
      background: #6e151b;
    }

    /* Colophon / Thanksgiving Box */
    .colophon-block {
      text-align: center;
      margin: 1rem auto 0.5rem auto;
      padding: 0.75rem 1rem;
      border-top: 1.5px solid rgba(140, 29, 36, 0.3);
      max-width: 450px;
      page-break-inside: avoid;
    }
    .colophon-title {
      font-family: 'Cinzel', serif;
      font-size: 9pt;
      font-weight: 700;
      letter-spacing: 0.2em;
      color: #8C1D24;
      margin: 0.35rem 0 0.25rem 0;
    }
    .colophon-text {
      font-family: 'EB Garamond', serif;
      font-style: italic;
      font-size: 9.5pt;
      color: #3C3632;
      margin-bottom: 0.5rem;
    }
    .colophon-motto {
      font-family: 'Cinzel', serif;
      font-size: 7.5pt;
      letter-spacing: 0.2em;
      color: #C5A059;
    }
  </style>
  <script>
    window.addEventListener('DOMContentLoaded', () => {
      const triggerPrint = () => {
        ${autoPrint ? 'setTimeout(() => window.print(), 350);' : ''}
      };
      if (document.fonts && document.fonts.ready) {
        document.fonts.ready.then(triggerPrint);
      } else {
        setTimeout(triggerPrint, 800);
      }
    });
  </script>
</head>
<body>

  <!-- BROWSER PRINT FLOATING CONTROLS -->
  <div class="print-floating-bar no-print">
    <button class="print-btn" onclick="window.print()">
      <span>🖨️</span>
      <span>Print / Save as PDF</span>
    </button>
  </div>

  <!-- TITLE PAGE / FRONTISPIECE -->
  <div class="title-page">
    <div class="book-border"></div>
    <div class="book-border-inner"></div>
    <div class="corner-cross corner-tl">☩</div>
    <div class="corner-cross corner-tr">☩</div>
    <div class="corner-cross corner-bl">☩</div>
    <div class="corner-cross corner-br">☩</div>

    <div>
      <div class="title-head-ornament">☦ ☩ ☦</div>
      <h1 class="title-main">Sancta Missa</h1>
      <div class="title-sub">Tridentina · Ritus Romanus 1962</div>
      <div class="title-divider"></div>
      <div style="font-family: 'Cinzel', serif; font-size: 8.5pt; letter-spacing: 0.18em; color: #3C3632; text-transform: uppercase;">
        Traditional Latin Mass according to the 1962 Roman Missal
      </div>
    </div>

    <div class="feast-container">
      <div class="feast-label">Liturgy of the Day · ${dateStr}</div>
      <div class="feast-title">${feast}</div>
      ${commemoration ? `<div class="feast-commemoration">Commemoratio: ${commemoration}</div>` : ''}
      <div style="margin-top: 0.65rem; font-family: 'Cinzel', serif; font-size: 7.75pt; color: #8C1D24; letter-spacing: 0.15em;">
        ${rubricsVersion}
      </div>
    </div>

    <div>
      <div style="font-size: 22pt; color: #C5A059; margin-bottom: 0.6rem;">✠</div>
      <div class="title-footer">
        Ad Maiorem Dei Gloriam · Divinum Officium
      </div>
    </div>
  </div>

  <!-- MAIN LITURGY -->
  <div class="content-container">
    <div class="table-column-headers">
      <div class="table-col-hdr">✠ TEXTUS LATINUS ✠</div>
      <div class="table-col-hdr">✠ ENGLISH TRANSLATION & RUBRICS ✠</div>
    </div>

    ${bodyHtml}
  </div>

</body>
</html>
`;
}
