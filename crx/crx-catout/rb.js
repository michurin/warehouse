// /radarbox
(() => {
  setInterval(() => {
    document.querySelectorAll('.map-ad').forEach(x => x.remove());
    document.querySelectorAll('.media-container').forEach(x => x.parentElement.remove());
  }, 2000);
})();
