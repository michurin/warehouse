(() => {
  setInterval(() => {
    document.querySelectorAll('#banner').forEach(x => x.remove());
  }, 2000);
})();
