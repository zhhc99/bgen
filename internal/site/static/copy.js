(function () {
  var COPY_SVG = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>';
  var CHECK_SVG = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>';

  function showCheck(btn) {
    btn.innerHTML = CHECK_SVG;
    btn.classList.add('copied');
    setTimeout(function () {
      btn.innerHTML = COPY_SVG;
      btn.classList.remove('copied');
    }, 2000);
  }

  function copyText(text, btn) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text).then(function () { showCheck(btn); });
      return;
    }
    // fallback for non-HTTPS environments
    var ta = document.createElement('textarea');
    ta.value = text;
    ta.style.cssText = 'position:fixed;opacity:0;pointer-events:none';
    document.body.appendChild(ta);
    ta.focus();
    ta.select();
    try { document.execCommand('copy'); showCheck(btn); } catch (e) {}
    document.body.removeChild(ta);
  }

  document.querySelectorAll('.copy-btn').forEach(function (btn) {
    btn.addEventListener('click', function () {
      var pre = btn.closest('pre');
      var code = pre && pre.querySelector('code');
      copyText((code || pre).innerText, btn);
    });
  });
})();
