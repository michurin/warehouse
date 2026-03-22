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
  return '#' + (Math.round((Math.random() + 1) * 4096)).toString(16).slice(-3)
}

const bodyElement = document.body
const rootElement = document.getElementById('board')
const statusElement = document.getElementById('status')
const colorElement = document.getElementById('color')
const nameElement = document.getElementById('name')
const inputElement = document.getElementById('input')
let color = localStorage.getItem('color') || randomColor()

colorElement.style.color = color
colorElement.onclick = () => {
  color = randomColor()
  colorElement.style.color = color
}
nameElement.value = localStorage.getItem('name') || 'me'
inputElement.focus()

function bar(text) {
  statusElement.textContent = text
}

bar('loading...')

// --- sending

inputElement.onkeyup = async (e) => {
  if (e.which === 10 || e.which === 13) {
    localStorage.setItem('name', nameElement.value)
    localStorage.setItem('color', color)
    const msg = inputElement.value.replaceAll(/[\n\r\t\s]/g, ' ')
    if (msg === '') {
      return
    }
    await fetch('/send', {
      method: 'POST',
      body: color + '|' + nameElement.value.replaceAll(/[\n\r\t\s]/g, ' ') + ': ' + inputElement.value.replaceAll(/[\n\r\t\s]/g, ' '),
    })
    inputElement.value = ''
    inputElement.focus()
  }
}

// --- fetching

const evtSource = new EventSource('/fetch')

evtSource.onmessage = (e) => {
  const a = e.data.split(/[\n\r]+/)
  a.reverse()
  a.forEach((text) => {
    while (rootElement.children.length > 10) {
      rootElement.firstChild.remove()
    }
    const newElement = document.createElement('div')
    newElement.textContent = text.slice(5)
    newElement.style.color = text.slice(0, 4)
    rootElement.append(newElement)
  })
}

evtSource.onerror = () => {
  bar('connection lost, wait a moment...')
}

evtSource.onopen = () => {
  bar('online')
}
