interface EventContext {
  request: Request;
  params: Record<string, string | string[]>;
}

export async function onRequestGet(context: EventContext): Promise<Response> {
  const url = new URL(context.request.url);
  const dateParam = url.searchParams.get('date') || 'today';
  const acceptHeader = context.request.headers.get('Accept') || '';

  // If client specifically requests JSON (e.g. API probe)
  if (acceptHeader.includes('application/json')) {
    return new Response(JSON.stringify({
      status: 'redirect_to_print',
      message: 'Headless Chrome PDF rendering is not hosted on edge workers. Use the print-ready endpoint to save high-fidelity vector PDF.',
      printUrl: `/api/missa-html?date=${encodeURIComponent(dateParam)}&print=true`
    }), {
      status: 200,
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*'
      }
    });
  }

  // Redirect browser navigation directly to print-ready layout
  const targetUrl = new URL('/api/missa-html', url.origin);
  targetUrl.searchParams.set('date', dateParam);
  targetUrl.searchParams.set('print', 'true');

  return Response.redirect(targetUrl.toString(), 302);
}
