const params = new URLSearchParams(window.location.search)
const room = params.get('back')
if (room) {
  const r = room.replaceAll(/[^0-9a-zA-Z_-]+/g, '')
  document.getElementsByTagName('b')[0].innerText = r
  document.getElementsByTagName('a')[0].href = '/' + r
}

