/*global fetch, EventSource, requestAnimationFrame, document:readable, setTimeout:readable, localStorage:readable*/
/*eslint indent: ["error", 2]*/
/*eslint eqeqeq: ["error", "always"]*/
/*eslint prefer-const: "error"*/
/*eslint no-var: "error"*/
/*eslint no-undef: "error"*/
/*eslint one-var: ["error", "never"]*/
/*eslint semi: ["error", "never"]*/
/*eslint quotes: ["error", "single"]*/
/*eslint prefer-arrow-callback: "error"*/
/*eslint arrow-body-style: "error"*/

// TODO languages: navigator.languages

function pad(x) { // naive
  return (x < 10 ? '0' : '') + x
}

function timeFormat(ts) {
  const d = new Date(ts)
  return pad(d.getHours()) + ':' + pad(d.getMinutes())
}

const boardElement = document.getElementById('board')
const statusElement = document.getElementById('status')
const colorElement = document.getElementById('color')
const nameElement = document.getElementById('name')
const inputElement = document.getElementById('input')
const sendElement = document.getElementById('send')
const lockElement = document.getElementById('lock')
const usersElement = document.getElementById('users')

const appState = {
  room: '',
  user: '',
  name: '', // TODO store in element?
  color: '', // TODO store in element?
  users: [],
  locked: false,
}

function initAppState() {
  appState.room = location.pathname.replaceAll(/[^0-9a-zA-Z_-]+/g, '') || 'main'
  appState.user = localStorage.getItem('user')
  appState.name = localStorage.getItem('name')
  appState.color = localStorage.getItem('color')
  if (!appState.user) {
    appState.user = Date.now().toString(36) + '-' + Math.random().toString(36).substring(2)
    localStorage.setItem('user', appState.user)
  }
  if (!appState.name) {
    appState.name = 'u' + ((Math.random() + 1) * 100).toString(10).substring(1, 3)
    localStorage.setItem('name', appState.name)
  }
  if (!/^#[0-9a-fA-F]{3,6}$/.test(appState.color)) {
    appState.color = '#' + ((Math.random() + 1) * 16777216).toString(16).substring(1, 7)
    localStorage.setItem('color', appState.color)
  }
  nameElement.value = appState.name
  colorElement.value = appState.color
}

function setLock(s) {
  appState.locked = s
  lockElement.textContent = s ? '🔐' : '🔓'
}

function setUsers(uu) {
  usersElement.innerHTML = ''
  uu.sort((a, b) => a.name.localeCompare(b.name) || a.color.localeCompare(b.color))
  appState.users = uu
  uu.forEach((u) => {
    const e = document.createElement('div')
    e.textContent = u.name
    e.style.color = u.color
    usersElement.append(e)
  })
}

async function send() {
  localStorage.setItem('name', nameElement.value) // TODO set default if empty
  localStorage.setItem('color', colorElement.value) // TODO validate, set default
  const msg = inputElement.value.replaceAll(/\p{Cc}+/gu, ' ')
  if (msg === '') {
    return
  }
  await fetch('/pub', {
    method: 'POST',
    body: JSON.stringify({
      room: appState.room,
      user: appState.user,
      color: colorElement.value.replaceAll(/\p{Cc}+/gu, ' '),
      name: nameElement.value.replaceAll(/\p{Cc}+/gu, ' '),
      message: msg,
    })
  })
  inputElement.value = ''
  inputElement.focus()
}

function eventMessage(e) {
  const a = e.data.split(/[\n\r]+/)
  a.reverse()
  a.forEach((bytes) => {
    const dto = JSON.parse(bytes)
    const m = dto.message
    if (m) {
      if (m.name === '#CONTROL') {
        location.href = '/fin.html?back=' + encodeURIComponent(appState.room)
        return
      }
      while (boardElement.children.length > 1000) {
        boardElement.firstChild.remove()
      }
      const eDiv = document.createElement('div')
      const eTS = document.createElement('code')
      const eB = document.createElement('b')
      const eSpan = document.createElement('span')
      eTS.textContent = timeFormat(m.ts) + ' '
      eB.textContent = m.name + ': '
      eSpan.textContent = m.message
      eDiv.append(eTS, eB, eSpan)
      eDiv.style.color = m.color
      boardElement.append(eDiv)
    } else {
      const eDiv = document.createElement('div') // TODO for debugging only!
      eDiv.textContent = JSON.stringify(dto)
      boardElement.append(eDiv)
    }
    console.log('dto', dto)
    if (dto.users) {
      setUsers(dto.users)
    }
    if (dto.locked !== undefined) {
      setLock(dto.locked)
    }
  })
}

function bar(text, title) {
  statusElement.textContent = text
  statusElement.title = title
}

function eventError() {
  bar('❌', 'offline')
}

function eventOpen() {
  bar('✅', 'online')
}

function toggleLock() {
  console.log('toggle lock')
  fetch('/lock', {
    method: 'POST',
    body: JSON.stringify({
      room: appState.room,
      user: appState.user,
      lock: !appState.locked,
    })
  })
}

function inputKeyup(e) {
  if (e.which === 10 || e.which === 13) { // it won't work on Androd Chrome
    send()
  }
}

function initApp() {
  const queryString = new URLSearchParams({ room: appState.room, user: appState.user }).toString()
  const evtSource = new EventSource('/fetch?' + queryString)
  evtSource.onmessage = eventMessage
  evtSource.onerror = eventError
  evtSource.onopen = eventOpen
  lockElement.onclick = toggleLock
  sendElement.onclick = send
  sendElement.ontouchstart = send // android
  sendElement.onmousedown = send // android with chrome bug
  inputElement.onkeyup = inputKeyup
  inputElement.focus()
}

(async function() {
  initAppState()
  const resp = await fetch('/enter', {
    method: 'POST',
    body: JSON.stringify({
      room: appState.room,
      user: appState.user,
      name: appState.name,
      color: appState.color,
    })
  })
  // TODO check resp.ok
  const data = await resp.json()
  setLock(data.locked)
  setUsers(data.users)
  initApp()
})()
