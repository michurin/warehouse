// /pogoda
(() => {
  setInterval(() => {
    document.querySelectorAll('.content__right').forEach(e => e.remove());
    document.querySelectorAll('.content__adv').forEach(e => e.remove());
    document.querySelectorAll('.card_with-horizontal-extension').forEach(e => e.remove());
  }, 2000);
})();
