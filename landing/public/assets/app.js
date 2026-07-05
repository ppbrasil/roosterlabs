(function () {
  var form = document.getElementById('lead-form');
  if (!form) return;

  var steps = Array.prototype.slice.call(form.querySelectorAll('.q-step'));
  var bars = Array.prototype.slice.call(form.querySelectorAll('.progress i'));
  var backBtn = form.querySelector('.back-link');
  var errBox = form.querySelector('.form-error');
  var current = 0;
  var answers = {};

  function show(i) {
    current = i;
    steps.forEach(function (s, idx) { s.style.display = idx === i ? 'block' : 'none'; });
    bars.forEach(function (b, idx) { b.classList.toggle('done', idx <= i); });
    backBtn.style.visibility = i === 0 ? 'hidden' : 'visible';
    errBox.style.display = 'none';
  }

  function next() {
    if (current < steps.length - 1) show(current + 1);
  }

  steps.forEach(function (step) {
    var radios = step.querySelectorAll('input[type="radio"]');
    radios.forEach(function (r) {
      r.addEventListener('change', function () {
        answers[r.name] = r.value;
        var other = step.querySelector('input[type="text"][data-other="' + r.name + '"]');
        if (r.value === 'other' && other) { other.focus(); return; }
        setTimeout(next, 220);
      });
    });
    var others = step.querySelectorAll('input[type="text"][data-other]');
    others.forEach(function (o) {
      o.addEventListener('keydown', function (e) {
        if (e.key === 'Enter') {
          e.preventDefault();
          answers[o.dataset.other] = 'other: ' + o.value.trim();
          next();
        }
      });
      o.addEventListener('blur', function () {
        if (o.value.trim()) answers[o.dataset.other] = 'other: ' + o.value.trim();
      });
    });
  });

  backBtn.addEventListener('click', function (e) {
    e.preventDefault();
    if (current > 0) show(current - 1);
  });

  form.addEventListener('submit', function (e) {
    e.preventDefault();
    var email = form.querySelector('#email');
    var linkedin = form.querySelector('#linkedin');
    if (!email.value || !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(email.value)) {
      errBox.textContent = form.dataset.errEmail;
      errBox.style.display = 'block';
      return;
    }
    var params = new URLSearchParams(window.location.search);
    var utm = {};
    ['utm_source', 'utm_medium', 'utm_campaign', 'utm_content', 'utm_term'].forEach(function (k) {
      if (params.get(k)) utm[k] = params.get(k);
    });
    var payload = {
      email: email.value.trim(),
      linkedin_url: linkedin.value.trim(),
      profile: answers.profile || null,
      goal: answers.goal || null,
      maturity: answers.maturity || null,
      challenge: answers.challenge || null,
      lang: document.documentElement.lang,
      utm: JSON.stringify(utm),
      website: form.querySelector('.hp input').value
    };
    var btn = form.querySelector('.btn');
    btn.disabled = true;
    btn.textContent = '…';
    fetch('/api/subscribe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(function (r) {
      if (!r.ok) throw new Error('bad status');
      form.style.display = 'none';
      document.getElementById('form-success').style.display = 'block';
    }).catch(function () {
      btn.disabled = false;
      btn.textContent = form.dataset.ctaLabel;
      errBox.textContent = form.dataset.errGeneric;
      errBox.style.display = 'block';
    });
  });

  show(0);
})();
