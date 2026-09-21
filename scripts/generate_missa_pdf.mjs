import fs from 'fs';
import path from 'path';
import { spawn } from 'child_process';
import { fileURLToPath } from 'url';
import {
  parseDateArg,
  fetchMissa,
  cleanCell,
  parseMissaHtml,
  buildHtmlDocument
} from '../frontend/src/utils/missa.ts';

// Re-export shared core logic
export { parseDateArg, fetchMissa, cleanCell, parseMissaHtml, buildHtmlDocument };

export function findChromePath() {
  const candidates = [
    process.env.CHROME_BIN,
    process.env.PUPPETEER_EXECUTABLE_PATH,
    '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
    '/Applications/Chromium.app/Contents/MacOS/Chromium',
    '/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary',
    '/usr/bin/google-chrome',
    '/usr/bin/google-chrome-stable',
    '/usr/bin/chromium',
    '/usr/bin/chromium-browser',
    'C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe',
    'C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe'
  ];

  for (const candidate of candidates) {
    if (candidate && fs.existsSync(candidate)) {
      return candidate;
    }
  }
  return null;
}

export async function renderPdfWithChrome(htmlPath, pdfPath) {
  const chromePath = findChromePath();
  if (!chromePath) {
    throw new Error('Google Chrome or Chromium binary was not found on this system.');
  }

  console.log(`[Chrome Headless] Using browser at ${chromePath}`);
  console.log(`[Chrome Headless] Rendering ${htmlPath} to ${pdfPath}...`);
  if (fs.existsSync(pdfPath)) {
    fs.unlinkSync(pdfPath);
  }

  const tempProfileDir = path.join('/tmp', `chrome-pdf-${Date.now()}`);

  const child = spawn(chromePath, [
    '--headless=new',
    '--disable-gpu',
    '--no-sandbox',
    '--no-first-run',
    '--no-default-browser-check',
    '--no-pdf-header-footer',
    `--user-data-dir=${tempProfileDir}`,
    `--print-to-pdf=${pdfPath}`,
    htmlPath
  ], { stdio: 'ignore' });

  // Poll for PDF file creation
  const maxAttempts = 60;
  for (let i = 0; i < maxAttempts; i++) {
    await new Promise(res => setTimeout(res, 250));
    if (fs.existsSync(pdfPath)) {
      const stats = fs.statSync(pdfPath);
      if (stats.size > 2000) {
        await new Promise(res => setTimeout(res, 500));
        child.kill('SIGKILL');
        try {
          fs.rmSync(tempProfileDir, { recursive: true, force: true });
        } catch (_) {}
        console.log(`[Chrome Headless] Successfully created PDF (${stats.size} bytes).`);
        return true;
      }
    }
  }

  child.kill('SIGKILL');
  throw new Error('PDF generation timed out after 15 seconds.');
}

async function main() {
  const args = process.argv.slice(2);
  const inputDate = args[0] || 'today';
  const { divinumDate, isoDate } = parseDateArg(inputDate);

  console.log(`=======================================================`);
  console.log(`Traditional Latin Mass (Divinum Officium) PDF Generator`);
  console.log(`Target Date: ${divinumDate} (ISO: ${isoDate})`);
  console.log(`Rubrics: English Only`);
  console.log(`=======================================================`);

  let htmlContent = '';
  try {
    htmlContent = await fetchMissa(divinumDate);
  } catch (err) {
    console.warn(`[Warning] Remote fetch failed: ${err.message}. Checking local sample cache...`);
    if (fs.existsSync('/tmp/missa_sample.html')) {
      htmlContent = fs.readFileSync('/tmp/missa_sample.html', 'utf8');
      console.log(`[Cache] Loaded sample cache.`);
    } else {
      throw err;
    }
  }

  console.log(`[Parser] Parsing Divinum Officium Mass structure...`);
  const parsed = parseMissaHtml(htmlContent);
  console.log(`Feast: ${parsed.feast}`);
  if (parsed.commemoration) console.log(`Commemoration: ${parsed.commemoration}`);
  console.log(`Total sections/rows: ${parsed.items.length}`);

  const formattedDate = new Date(`${isoDate}T12:00:00`).toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });

  const fullHtml = buildHtmlDocument({
    feast: parsed.feast,
    commemoration: parsed.commemoration,
    rubricsVersion: parsed.rubricsVersion,
    items: parsed.items,
    dateStr: formattedDate
  });

  const tempHtmlPath = path.join('/tmp', `missa_${isoDate}.html`);
  fs.writeFileSync(tempHtmlPath, fullHtml, 'utf8');
  console.log(`[HTML] Written temporary print layout to ${tempHtmlPath}`);

  // Determine output paths
  const outputFileName = `Traditional_Latin_Mass_${isoDate}.pdf`;
  const rootPdfPath = path.join(process.cwd(), outputFileName);

  await renderPdfWithChrome(tempHtmlPath, rootPdfPath);

  console.log(`=======================================================`);
  console.log(`✓ PDF Generated Successfully!`);
  console.log(`File: ${rootPdfPath}`);
  console.log(`Size: ${(fs.statSync(rootPdfPath).size / 1024).toFixed(1)} KB`);
  console.log(`=======================================================`);
}

// Only run main if invoked directly
if (process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1]) {
  main().catch(err => {
    console.error('Fatal Error:', err);
    process.exit(1);
  });
}
