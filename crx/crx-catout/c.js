(() => {
  setInterval(() => {
    ['blocked-results-banner', 'rca', 'lig_reverso_smartbox_article_tc'].forEach((x) => {
      const e = document.getElementById(x);
      if (e) {
        e.remove();
      }
    });
    document.querySelectorAll('.blocked').forEach(x => x.classList.remove('blocked'));
  }, 2000);
})();
