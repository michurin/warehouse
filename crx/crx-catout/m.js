(() => {

  return; // --------------------------------

  setInterval(() => {
    // letter-like
    let cc = document.querySelector('[data-key="box=toolbar-box"]');
    if (cc) {
      cc = cc.parentElement.children;
      if (cc.length === 5) {
        cc[3].remove();
      }
    }
    // promo
    document.querySelectorAll('[data-key="view=mail-pro-left-column-button"]').forEach(e => e.remove());
    document.querySelectorAll('.PSHeader-Pro').forEach(e => e.remove());
    // left banners first generation
    document.querySelectorAll('.b-banner').forEach(e => e.remove());
    // left banners second and further generations
    cc = document.querySelectorAll('.mail-NestedList-Setup');
    if (cc.length === 3) {
      const x = document.querySelectorAll('.mail-NestedList-Setup')[2].nextElementSibling;
      if (x) {
        x.remove();
      }
    }
    // top ad line
    cc = document.querySelectorAll('.mail-Layout-Main');
    if (cc.length === 1) {
      cc = cc[0].children;
      if (cc.length === 3) {
        cc[1].remove();
      }
    }
    // remove annoying and useless garbage
    document.querySelectorAll('.PSHeaderLogo360').forEach(e => e.remove());
    document.querySelectorAll('.PSHeader-Center').forEach(e => e.style.opacity = 0);
    // annoying banner YA360
    document.querySelectorAll('.WithSlidingSearch').forEach((e, n) => { if (n === 0) { e.style.opacity = 0; } });
    // the ad line right below search and above controls like reply etc
    document.querySelectorAll('.mail-DirectLineContainer').forEach(e => {
      e.style.position = 'absolute';
      e.offsetTop = '-1000px';
    });
  }, 2000);
})();
