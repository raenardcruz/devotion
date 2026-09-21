import {
  parseDateArg,
  fetchMissa,
  parseMissaHtml,
  buildHtmlDocument
} from '../../src/utils/missa';

interface EventContext {
  request: Request;
  params: Record<string, string | string[]>;
}

export async function onRequestGet(context: EventContext): Promise<Response> {
  const url = new URL(context.request.url);
  const dateParam = url.searchParams.get('date') || 'today';
  const autoPrint = url.searchParams.get('print') === 'true' || url.searchParams.get('download') === 'true';

  try {
    const { divinumDate, isoDate } = parseDateArg(dateParam);
    const rawMissaHtml = await fetchMissa(divinumDate);
    const parsed = parseMissaHtml(rawMissaHtml);

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

    return new Response(fullHtml, {
      status: 200,
      headers: {
        'Content-Type': 'text/html; charset=utf-8',
        'Access-Control-Allow-Origin': '*',
        'Cache-Control': 'public, max-age=86400, s-maxage=86400'
      }
    });
  } catch (err: any) {
    console.error('[missa-html error]:', err);
    const errorHtml = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Liturgy Error - Sancta Missa</title>
  <style>
    body { font-family: serif; max-width: 600px; margin: 4rem auto; padding: 2rem; text-align: center; color: #1A1715; }
    h1 { color: #8C1D24; font-size: 1.5rem; }
    p { font-size: 0.95rem; color: #555; }
    a { display: inline-block; margin-top: 1.5rem; color: #8C1D24; font-weight: bold; text-decoration: none; border: 1px solid #C5A059; padding: 0.5rem 1rem; border-radius: 8px; }
  </style>
</head>
<body>
  <h1>☦ Sancta Missa Retrieval Error ☦</h1>
  <p>Unable to retrieve liturgical text for date: <strong>${dateParam}</strong></p>
  <p>${err.message || 'Unknown error occurred while connecting to Divinum Officium.'}</p>
  <a href="javascript:history.back()">← Return to Daily Readings</a>
</body>
</html>`;

    return new Response(errorHtml, {
      status: 502,
      headers: {
        'Content-Type': 'text/html; charset=utf-8',
        'Access-Control-Allow-Origin': '*'
      }
    });
  }
}
