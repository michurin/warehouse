function log(text, color) {
  const e = document.createElement('div');
  e.innerText = text;
  if (color) {
    e.style.color = color;
  }
  document.getElementById('log').prepend(e);
}

function safe(name, f) {
  return async function () {
    log(name, '#080');
    try {
      await f();
      log(`${name} ok. done.`, '#080');
    } catch (e) {
      log(`${name} ${e}`, '#800');
    }
  };
}

document.getElementById('export').onclick = safe('EXPORT', async () => {
  log('feaching data...');
  const data = localStorage.days || '{}';
  log(data);
  document.getElementById('json').innerText = data;
});

document.getElementById('import').onclick = safe('IMPORT', async () => {
  log('feaching original data...');
  const data = localStorage.days || '{}';
  log(data);
  log('feaching update from clipboard...');
  const update = document.getElementById('json').innerText;
  log(update);
  log('checking...');
  const s = JSON.parse(update);
  log('json is valid');
  localStorage.setItem('days', JSON.stringify(s));
  log('updated. refrash your calendar page');
});

document.getElementById('merge').onclick = safe('MERGE', async () => {
  log('not implemented');
});

document.getElementById('merge_dry').onclick = safe('MERGE DRY RUN', async () => {
  log('not implemented');
});

document.getElementById('undo').onclick = safe('UNDO', async () => {
  log('not implemented');
});
