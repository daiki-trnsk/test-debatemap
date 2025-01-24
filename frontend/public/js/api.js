const API_HOST = 'http://localhost:8000'

// const logout = () => {
//   localStorage.removeItem('token')
//   location.href = '/login.html'
// }

const handleLoginError = () => {
  alert('ログイン情報が確認できませんでした、ログインページに移動します')
  logout()
}

const handleForbiddenError = () => {
  alert('他のユーザーのタスクは閲覧できません')
  throw new Error('他のユーザーのタスクは閲覧できません')
}

const handleOtherError = () => {
  alert('予期せむエラーが発生しました')
  throw new Error('予期せむエラーが発生しました')
}

/**
 * ユーザー新規登録API
 */
const signUpApi = (data) => {
  const url = `${API_HOST}/user`
  return fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  }).then((res) => {
    if (res.ok) {
      return res.json()
    } else if (res.status === 400) {
      console.error(res)
      throw new Error('入力されたメールアドレスは既に登録されています')
    } else {
      console.error(res)
      handleOtherError()
    }
  })
}

/**
 * ログインAPI
 */
const loginApi = (email, password) => {
  const url = `${API_HOST}/token`
  return fetch(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    // loginの場合のみ、bodyは特別
    body: `username=${email}&password=${password}`,
  })
    .then((res) => {
      if (res.ok) {
        return res.json()
      } else if (res.status === 401) {
        console.error(res)
        throw new Error('メールアドレスまたはパスワードが間違っています')
      } else {
        console.error(res)
        handleOtherError()
      }
    })
    .then((data) => {
      localStorage.setItem('token', data.access_token)
      return data
    })
}

/**
 * ログインユーザー情報取得するAPI
 */
const getMeApi = () => {
  const url = `${API_HOST}/me`
  return fetch(url, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  }).then((res) => {
    if (res.ok) {
      return res.json()
    } else if (res.status === 401) {
      handleLoginError()
    } else {
      console.error(res)
      handleOtherError()
    }
  })
}