'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import '../../styles/register.css'

export default function SignupPage() {
  const [form, setForm] = useState({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dob: '',
    avatar: '',
    nickname: '',
    aboutMe: ''
  })

  const [error, seterror] = useState('')
  const router = useRouter()

  const handleChange = e => {
    setForm({
      ...form,
      [e.target.name]: e.target.name === "dob"
        ? Math.floor(new Date(e.target.value).getTime() / 1000)
        : e.target.value

    })
  }

  const handleFileChange = e => {
    setForm({ ...form, avatar: e.target.files[0] })
  }

  const handleSignup = async e => {
    e.preventDefault()

    const formData = new FormData()
    formData.append("email", form.email)
    formData.append("password", form.password)
    formData.append("firstName", form.firstName)
    formData.append("lastName", form.lastName)
    formData.append("dob", form.dob)
    formData.append("nickname", form.nickname)
    formData.append("aboutMe", form.aboutMe)
    formData.append("avatar", form.avatar)
    formData.append("privacy", form.privacy || "public")

    try {
      const response = await fetch("http://localhost:8080/api/register", {
        method: "POST",
        body: formData,
      })
  

      const res = await response.json()
      if  (res.error){

        seterror(res.error)
        return 
      } 

      router.push('/login')

    } catch (error) {
    }
  }


  return (
    <div className="containerregister">
      <div className="bg"></div>
      <h1 className="titleregister">Create Account</h1>

      {error && <p className="error">{error}</p>}

      <form onSubmit={handleSignup} className="formregister">
        <div className="row">
          <input type="text" name="firstName" placeholder="First Name" required onChange={handleChange} />
          <input type="text" name="lastName" placeholder="Last Name" required onChange={handleChange} />
        </div>

        <input type="email" name="email" placeholder="Email" required onChange={handleChange} />
        <input type="password" name="password" placeholder="Password" required onChange={handleChange} />
        <input type="date" name="dob" required onChange={handleChange} />

        <input type="text" name="nickname" placeholder="Nickname (optional)" onChange={handleChange} />
        <input type="file" name="avatar" accept="image/*" onChange={handleFileChange} />
        <textarea name="aboutMe" placeholder="About me (optional)" onChange={handleChange}></textarea>

        <label>
          <input
            type="radio"
            name="privacy"
            value="private"
            onChange={handleChange}
            required
          />
          Private
        </label>

        <label>
          <input
            type="radio"
            name="privacy"
            value="public"
            onChange={handleChange}
          />
          Public
        </label>


        <button type="submit" className="button">Sign Up</button>
      </form>

      <p className="footer">
        Already have an account? <Link href="/login">Login</Link>
      </p>
    </div>
  )
}
