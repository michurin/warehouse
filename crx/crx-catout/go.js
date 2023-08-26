// https://go.dev/play/
(() => {
  var rn;
  const f = () => {
    document.querySelectorAll('#editor-description').forEach(x => { x.style.fontSize = '5px' });
    document.querySelectorAll('.Playground-inputContainer').forEach(x => {
      x.style.height = '60vh';
      x.scrollIntoView();
    });
    document.querySelectorAll('.Cookie-notice').forEach(x => { x.style.display = 'none' });
    document.querySelectorAll('.Playground-title').forEach(x => { x.style.display = 'none' });
    document.querySelectorAll('#run').forEach(x => {
      if (rn) {
        return;
      }
      x.style.position = 'fixed';
      x.style.top = 0;
      x.style.right = 0;
      x.style.zIndex = 999;
      x.innerText += " (ctrl-space)";
      rn = x;
    });
    document.querySelectorAll('#code').forEach(x => {
      if (x.value.substr(0, 2) === '//') {
        x.value = x.value.replaceAll(/^\/\/.*$/gm, '').replace(/^\s+/, '');
        x.onkeyup = (e) => {
          if (e.keyCode == 32 && (e.ctrlKey || e.metaKey)) {
            rn.click();
          }
        }
        x.focus();
      }
    });
  };
  setTimeout(f, 200);
  setTimeout(f, 500);
  setTimeout(f, 1000);
  setTimeout(f, 3000);
})();
