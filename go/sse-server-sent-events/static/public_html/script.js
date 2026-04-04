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

function randomColor() {
  return '#' + (Math.round((Math.random() + 1) * 16777216)).toString(16).slice(-6)
}

function pad(x) { // naive
  return (x < 10 ? '0' : '') + x
}

function timeFormat(ts) {
  const d = new Date(ts)
  return pad(d.getHours()) + ':' + pad(d.getMinutes())
}

const bodyElement = document.body
const boardElement = document.getElementById('board')
const statusElement = document.getElementById('status')
const colorElement = document.getElementById('color')
const nameElement = document.getElementById('name')
const inputElement = document.getElementById('input')
const sendElement = document.getElementById('send')

colorElement.value = localStorage.getItem('color') || randomColor() // TODO validate
nameElement.value = localStorage.getItem('name') || 'me'
inputElement.focus()

function bar(text, title) {
  statusElement.textContent = text
  statusElement.title = title
}

const roomID = 'main' // TODO: get from URL with fallback

const userID = (function() {
  let u = localStorage.getItem('user')
  if (!u) {
    u = Date.now().toString(36) + '-' + Math.random().toString(36).substring(2)
    localStorage.setItem('user', u)
  }
  return u
})()
console.log(userID)

const queryString = new URLSearchParams({ room: roomID, user: userID }).toString()

bar('loading...')

// --- sending

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
      room: roomID,
      user: userID,
      color: colorElement.value.replaceAll(/\p{Cc}+/gu, ' '),
      name: nameElement.value.replaceAll(/\p{Cc}+/gu, ' '),
      message: msg,
    })
  })
  inputElement.value = ''
  inputElement.focus()
}

inputElement.onkeyup = (e) => {
  if (e.which === 10 || e.which === 13) { // it won't work on Androd Chrome
    send()
  }
}

sendElement.onclick = send
sendElement.ontouchstart = send // android
sendElement.onmousedown = send // android with chrome bug

// --- fetching

const evtSource = new EventSource('/fetch?' + queryString)

evtSource.onmessage = (e) => {
  const a = e.data.split(/[\n\r]+/)
  a.reverse()
  a.forEach((bytes) => {
    const dto = JSON.parse(bytes)
    if (dto.message) {
      while (boardElement.children.length > 1000) {
        boardElement.firstChild.remove()
      }
      const m = dto.message
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
      const eDiv = document.createElement('div')
      eDiv.textContent = JSON.stringify(dto)
      boardElement.append(eDiv)
    }
  })
}

evtSource.onerror = () => {
  bar('❌', 'offline')
}

evtSource.onopen = () => {
  bar('✅', 'online')
}
