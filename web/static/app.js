(function () {
  function readDataset(name) {
    const body = document.body;
    if (!body || !body.dataset) {
      return "";
    }
    return (body.dataset[name] || "").trim();
  }

  // UTM real do visitante atual vem só da URL, nunca do HTML renderizado no
  // servidor (épico 002, T16): o CloudFront cacheia GET / ignorando query
  // string, então tanto o dataset do <body> quanto os hidden inputs do passo
  // 1 do form podem ter sido gravados para OUTRO visitante que bateu nesse
  // mesmo cache. location.search é a única fonte confiável no cliente.
  const params = new URLSearchParams(window.location.search);
  const utm = {
    utm_source: (params.get("utm_source") || "").trim(),
    utm_medium: (params.get("utm_medium") || "").trim(),
    utm_campaign: (params.get("utm_campaign") || "").trim(),
    utm_term: (params.get("utm_term") || "").trim(),
    utm_content: (params.get("utm_content") || "").trim()
  };

  // Sobrescreve incondicionalmente os hidden inputs do passo 1 — inclusive
  // com "", pra não deixar vazar UTM de outro visitante servido do cache.
  // Passos 2+ são renderizados pelo servidor a partir do que foi submetido no
  // passo 1 (ver utmFromForm em form_handler.go), então corrigir aqui já
  // basta pro funil inteiro.
  const form = document.querySelector("#lead-form form");
  if (form) {
    Object.keys(utm).forEach(function (key) {
      const input = form.querySelector('input[name="' + key + '"]');
      if (input) {
        input.value = utm[key];
      }
    });
  }

  const payload = {
    path: readDataset("path") || window.location.pathname,
    language: readDataset("language") || "pt-BR",
    utm_source: utm.utm_source,
    utm_medium: utm.utm_medium,
    utm_campaign: utm.utm_campaign,
    utm_term: utm.utm_term,
    utm_content: utm.utm_content
  };

  navigator.sendBeacon(
    "/event/view",
    new Blob([JSON.stringify(payload)], { type: "application/json" })
  );
})();
