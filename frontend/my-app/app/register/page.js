'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import styles from "./register.module.css"

export default function SignupPage() {
  const [form, setForm] = useState({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    dob: '',
    avatar: null,
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
    try {
      const response = await fetch('http://localhost:8080/api/register ', {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(form),
      })

      if (response.ok) {
        router.push('/login')
        
      } else {
        console.log(1);
        const data = await response.json()
        seterror(data)
      }
    } catch (error) {
      console.log(error);

      seterror("Network error!", error)
    }
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Create Account</h1>

      {error && <p className={styles.error}>{error}</p>}

      <form onSubmit={handleSignup} className={styles.form}>
        <div className={styles.row}>
          <input type="text" name="firstName" placeholder="First Name" required onChange={handleChange} />
          <input type="text" name="lastName" placeholder="Last Name" required onChange={handleChange} />
        </div>

        <input type="email" name="email" placeholder="Email" required onChange={handleChange} />
        <input type="password" name="password" placeholder="Password" required onChange={handleChange} />
        <input type="date" name="dob" required onChange={handleChange} />

        <input type="text" name="nickname" placeholder="Nickname (optional)" onChange={handleChange} />
        <input type="file" name="avatar" accept="image/*" onChange={handleFileChange} />
        <textarea name="aboutMe" placeholder="About me (optional)" onChange={handleChange}></textarea>

        <label className={styles.checkbox}>
          <input type="checkbox" required />
          I agree to the Terms & Privacy Policy
        </label>

        <button type="submit" className={styles.button}>Sign Up</button>
      </form>

      <p className={styles.footer}>
        Already have an account? <Link href="/login">Login</Link>
      </p>
    </div>
  )
}
