import { defineConfig, loadEnv, type Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import fs from 'fs'

function kindleHtmlPlugin(mode: string): Plugin {
  const env = loadEnv(mode, process.cwd(), '');
  const apiBaseUrl = env.VITE_API_BASE_URL || process.env.VITE_API_BASE_URL || 'https://devotionapi.raenardcruz.com';
  const apiToken = env.VITE_API_TOKEN || process.env.VITE_API_TOKEN || '';

  const injectEnv = (html: string) => {
    return html
      .replace(/__VITE_API_BASE_URL__/g, apiBaseUrl)
      .replace(/__VITE_API_TOKEN__/g, apiToken);
  };

  return {
    name: 'kindle-html-plugin',
    configureServer(server) {
      server.middlewares.use((req, _res, next) => {
        if (req.url) {
          const url = new URL(req.url, 'http://localhost');
          if (url.pathname === '/kindle' || url.pathname === '/kindle/') {
            req.url = '/kindle.html' + url.search;
          }
        }
        next();
      });
    },
    transformIndexHtml(html) {
      return injectEnv(html);
    }
  };
}

function missaPdfPlugin(): Plugin {
  return {
    name: 'missa-pdf-plugin',
    configureServer(server) {
      server.middlewares.use(async (req, res, next) => {
        const url = new URL(req.url || '', 'http://localhost');
        if (url.pathname === '/api/missa-pdf' || url.pathname === '/api/missa-html') {
          const dateParam = url.searchParams.get('date') || 'today';
          const autoPrint = url.searchParams.get('print') === 'true';

          try {
            // @ts-ignore
            const { parseDateArg, fetchMissa, parseMissaHtml, buildHtmlDocument } = await import('./src/utils/missa');
            const { divinumDate, isoDate } = parseDateArg(dateParam);
            
            let htmlContent = '';
            try {
              htmlContent = await fetchMissa(divinumDate);
            } catch (fetchErr) {
              if (fs.existsSync('/tmp/missa_sample.html')) {
                htmlContent = fs.readFileSync('/tmp/missa_sample.html', 'utf8');
              } else {
                throw fetchErr;
              }
            }

            const parsed = parseMissaHtml(htmlContent);
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
              dateStr: formattedDate,
              autoPrint
            });

            if (url.pathname === '/api/missa-html') {
              res.setHeader('Content-Type', 'text/html; charset=utf-8');
              res.end(fullHtml);
              return;
            }

            // PDF output: Try headless Chrome if present, otherwise redirect to print view
            try {
              // @ts-ignore
              const { renderPdfWithChrome } = await import('../scripts/generate_missa_pdf.mjs');
              const tempHtmlPath = path.join('/tmp', `missa_${isoDate}.html`);
              fs.writeFileSync(tempHtmlPath, fullHtml, 'utf8');

              const pdfFileName = `Traditional_Latin_Mass_${isoDate}.pdf`;
              const rootPdfPath = path.resolve(__dirname, '..', pdfFileName);
              await renderPdfWithChrome(tempHtmlPath, rootPdfPath);

              if (fs.existsSync(rootPdfPath)) {
                const fileData = fs.readFileSync(rootPdfPath);
                res.setHeader('Content-Type', 'application/pdf');
                res.setHeader('Content-Disposition', `attachment; filename="${pdfFileName}"`);
                res.end(fileData);
                return;
              }
            } catch (chromeErr: any) {
              console.warn('[Missa PDF Plugin] Headless Chrome unavailable, falling back to printable layout:', chromeErr.message);
            }

            // Fallback: Redirect to printable layout with print dialog ready
            res.writeHead(302, { Location: `/api/missa-html?date=${encodeURIComponent(dateParam)}&print=true` });
            res.end();
          } catch (err: any) {
            console.error('[Missa PDF Plugin Error]:', err);
            res.statusCode = 500;
            res.end(err.message || 'Failed to generate Latin Mass PDF');
          }
        } else {
          next();
        }
      });
    }
  };
}

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [vue(), kindleHtmlPlugin(mode), missaPdfPlugin()],
  build: {
    rollupOptions: {
      input: {
        main: path.resolve(__dirname, 'index.html'),
        kindle: path.resolve(__dirname, 'kindle.html'),
      }
    }
  },
  server: {
    proxy: {
      '/mass-readings': 'http://localhost:8788'
    }
  }
}))

