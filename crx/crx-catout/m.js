(() => {
  setInterval(() => {
    // letter-like
    var cc = document.querySelector('[data-key="box=toolbar-box"]').parentElement.children;
    if (cc.length === 5) {
      cc[3].remove();
    }
    // promo
    document.querySelectorAll('[data-key="view=mail-pro-left-column-button"]').forEach(e => e.remove());
    // left banners first generation
    document.querySelectorAll('.b-banner').forEach(e => e.remove());
    // left banners second and further generations
    cc = document.querySelectorAll('.mail-NestedList-Setup');
    if (cc.length === 3) {
      var x = document.querySelectorAll('.mail-NestedList-Setup')[2].nextElementSibling
      if (x) {
        x.remove();
      }
    }
  }, 2000);
})();