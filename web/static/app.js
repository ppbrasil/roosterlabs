(function () {
  function readDataset(name) {
    const body = document.body;
    if (!body || !body.dataset) {
      return "";
    }
    return (body.dataset[name] || "").trim();
  }

  const payload = {
    path: readDataset("path") || window.location.pathname,
    language: readDataset("language") || "pt-BR",
    utm_source: readDataset("utmSource"),
    utm_medium: readDataset("utmMedium"),
    utm_campaign: readDataset("utmCampaign"),
    utm_term: readDataset("utmTerm"),
    utm_content: readDataset("utmContent")
  };

  navigator.sendBeacon(
    "/event/view",
    new Blob([JSON.stringify(payload)], { type: "application/json" })
  );
})();
