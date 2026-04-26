'use strict'

// TODO languages: navigator.languages

function pad(x) { // naive
  return (x < 10 ? '0' : '') + x
}

function timeFormat(ts) {
  const d = new Date(ts)
  return pad(d.getHours()) + ':' + pad(d.getMinutes())
}

function isCanonicalName(x) {
  return typeof x == 'string' && /^[0-9A-Za-z_-]{1,32}$/.test(x)
}

function isValidColor(x) {
  return typeof x == 'string' && /^#[0-9A-Fa-f]{6}$/.test(x)
}

function noNulGuard(x) { // for debugging only
  if (!x) {
    throw new Error('not true')
  }
  return x
}

const eBoard = noNulGuard(document.getElementById('board'))
const eStatus = noNulGuard(document.getElementById('status'))
const eColorInput = noNulGuard(document.getElementById('color'))
const eNameInput = noNulGuard(document.getElementById('name'))
const eForm = noNulGuard(document.getElementById('form'))
const eMessageInput = noNulGuard(document.getElementById('input'))
const eLock = noNulGuard(document.getElementById('lock'))
const eUsers = noNulGuard(document.getElementById('users'))
const eShowUsers = noNulGuard(document.getElementById('show-users'))

const appState = {
  room: '',
  user: '',
  users: [],
  locked: false,
}

function initAppState() {
  var room = location.pathname.replace(/^\//, '')
  if (!isCanonicalName(room)) {
    location.href = '/main'
    return
  }
  appState.room = room
  appState.user = localStorage.getItem('user')
  if (!isCanonicalName(appState.user)) {
    appState.user = Date.now().toString(36) + '-' + Math.random().toString(36).substring(2)
    localStorage.setItem('user', appState.user)
  }
  var name = localStorage.getItem('name')
  if (!isCanonicalName(name)) {
    name = 'name' + ((Math.random() + 1) * 100).toString(10).substring(1, 3)
    localStorage.setItem('name', name)
  }
  eNameInput.value = name
  var color = localStorage.getItem('color')
  if (!isValidColor(color)) {
    color = '#' + ((Math.random() + 1) * 16777216).toString(16).substring(1, 7)
    localStorage.setItem('color', color)
  }
  eColorInput.value = color
}

function setLock(s) {
  if (s === undefined) {
    s = false
    eLock.disabled = true
  }
  appState.locked = s
  eLock.textContent = s ? '🔐' : '🔓'
}

function setUsers(uu) {
  eUsers.innerHTML = ''
  uu.sort((a, b) => a.name.localeCompare(b.name) || a.color.localeCompare(b.color))
  appState.users = uu
  uu.forEach((u) => {
    const e = document.createElement('div')
    e.textContent = u.name
    e.style.color = u.color
    eUsers.append(e)
  })
}

async function send() {
  localStorage.setItem('name', eNameInput.value) // TODO set default if empty
  localStorage.setItem('color', eColorInput.value) // TODO validate, set default
  const msg = eMessageInput.value.replaceAll(/\p{Cc}+/gu, ' ')
  if (msg === '') {
    return
  }
  await fetch('bin/pub', {
    method: 'POST',
    body: JSON.stringify({
      room: appState.room, // TODO check: if not -> fin.html
      user: appState.user, // TODO check: if not -> fix and focus name
      color: eColorInput.value.replaceAll(/\p{Cc}+/gu, ' '),
      name: eNameInput.value.replaceAll(/\p{Cc}+/gu, ' '),
      message: msg,
    }),
  })
  // TODO consider response, user can be kicked and publishing can be dropped
  eMessageInput.value = ''
  eMessageInput.focus()
}

function eventMessage(e) {
  const a = e.data.split(/[\n\r]+/)
  a.reverse()
  a.forEach((bytes) => {
    const dto = JSON.parse(bytes)
    const m = dto.message
    if (m) {
      if (m.name === '#CONTROL') {
        location.href = 's/fin.html?back=' + encodeURIComponent(appState.room) + '&reason=timeout'
        return
      }
      while (eBoard.children.length > 1000) {
        eBoard.firstChild.remove()
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
      eBoard.append(eDiv)
      eBoard.scrollTop = eBoard.scrollHeight
    }
    if (dto.users) {
      setUsers(dto.users)
    }
    if (dto.locked !== undefined) {
      setLock(dto.locked)
    }
  })
}

function bar(text, title) {
  eStatus.textContent = text
  eStatus.title = title
}

function eventError() {
  bar(appState.room + ': ❌', 'offline')
}

function eventOpen() {
  bar(appState.room + ': ✅', 'online')
}

function toggleLock() {
  console.log('toggle lock')
  fetch('bin/lock', {
    method: 'POST',
    body: JSON.stringify({
      room: appState.room,
      user: appState.user,
      lock: !appState.locked,
    }),
  })
}

function formSubmit(e) {
  e.preventDefault()
  send()
}

function inputKeyDown(e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    eForm.requestSubmit()
  }
}

function initApp() {
  const queryString = new URLSearchParams({ room: appState.room, user: appState.user }).toString()
  const evtSource = new EventSource('bin/fetch?' + queryString)
  evtSource.onmessage = eventMessage
  evtSource.onerror = eventError
  evtSource.onopen = eventOpen
  eLock.onclick = toggleLock
  eForm.onsubmit = formSubmit
  eMessageInput.onkeydown = inputKeyDown
  eMessageInput.focus()
  eShowUsers.onclick = () => { eUsers.style.display = eUsers.style.display === 'none' ? 'block' : 'none' }
  eUsers.onclick = () => { eUsers.style.display = 'none' }
}

(async function() { // TODO onpageshow? we have to do it on [back] action as well
  initAppState()
  const resp = await fetch('bin/enter', {
    method: 'POST',
    body: JSON.stringify({
      room: appState.room,
      user: appState.user,
      name: eNameInput.value,
      color: eColorInput.value,
    }),
  })
  // TODO check resp.ok
  const data = await resp.json()
  if (data.message) {
    location.href = 's/fin.html?back=main&reason=' + data.message.message
  }
  setLock(data.locked)
  setUsers(data.users)
  initApp()
})()
