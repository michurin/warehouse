// avito link in ms exchange
(() => {
  const F = () => {
    Array.from(document.querySelectorAll('*'))
      .map(x => (x.childNodes.length === 1 ? (x.childNodes[0].nodeType === 3 ? (x.innerText.match('^(http|avito.zoom.us/)') ? x : false) : false) : false))
      .filter(x => x)
      .forEach(x => {
        const c = x.childNodes[0];
        let t = x.innerText;
        if (!t.startsWith('http')) { // heuristic for naked domains
          t = `https://${t}`;
        }
        x.removeChild(c);
        const a = document.createElement('a');
        a.href = t;
        a.innerText = t.replace(/[^/]*:\/\//, '');
        a.target = '_blank';
        a.style.color = '#080';
        a.style.fontWeight = 'bold';
        x.append(a);
      });
  };
  setTimeout(F, 500);
  setTimeout(F, 1000);
  setInterval(F, 2000);
})();
