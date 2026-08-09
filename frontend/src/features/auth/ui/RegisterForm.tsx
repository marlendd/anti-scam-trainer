import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { ApiError } from '@/shared/api'

import { Input, PasswordInput } from '@/shared/ui/input'

import styles from './AuthForm.module.scss'
import { Logo } from '@/shared/ui/logo'
import { useRegister } from '@/features/auth'

export function RegisterForm() {
    const [confirmationError, setConfirmationError] = useState<string>()

    const register = useRegister()
    const navigate = useNavigate()

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        const formData = new FormData(event.currentTarget)

        const name = String(formData.get('name') ?? '')
        const email = String(formData.get('email') ?? '')
        const password = String(formData.get('password') ?? '')
        const passwordConfirmation = String(formData.get('passwordConfirmation') ?? '')

        if (password !== passwordConfirmation) {
            setConfirmationError('Пароли не совпадают')
            return
        }

        setConfirmationError(undefined)

        register.mutate(
            {
                name,
                email,
                password,
            },
            {
                onSuccess: () => {
                    navigate('/home', { replace: true })
                },
            },
        )
    }

    let errorMessage: string | undefined

    if (register.error instanceof ApiError) {
        if (register.error.status === 401) {
            errorMessage = 'Неверная электронная почта или пароль'
        } else {
            errorMessage = 'Не удалось войти. Попробуйте ещё раз'
        }
    } else if (register.isError) {
        errorMessage = 'Произошла ошибка. Попробуйте ещё раз'
    }

    return (
        <div className={styles.content}>
            <header className={styles.header}>
                <Link to="/home" className={styles.brand}>
                    <Logo />
                </Link>

                <p className={styles.description}>Создайте аккаунт для сохранения прогресса</p>
            </header>

            <form className={styles.form} onSubmit={handleSubmit}>
                <Input
                    label="Имя"
                    name="name"
                    placeholder="Введите имя"
                    autoComplete="name"
                    minLength={2}
                    required
                />

                <Input
                    label="Электронная почта"
                    name="email"
                    type="email"
                    placeholder="example@mail.ru"
                    autoComplete="email"
                    inputMode="email"
                    required
                />

                <PasswordInput
                    label="Пароль"
                    name="password"
                    placeholder="Не менее 8 символов"
                    autoComplete="new-password"
                    minLength={8}
                    required
                />

                <PasswordInput
                    label="Повторите пароль"
                    name="passwordConfirmation"
                    placeholder="Введите пароль ещё раз"
                    autoComplete="new-password"
                    minLength={8}
                    error={confirmationError}
                    onChange={() => {
                        if (confirmationError) {
                            setConfirmationError(undefined)
                        }
                    }}
                    required
                />

                <button className={styles.submit} type="submit">
                    Создать аккаунт
                </button>

                {errorMessage && <p className={styles.error}>{errorMessage}</p>}
            </form>

            <p className={styles.footer}>
                Уже есть аккаунт?{' '}
                <Link to="/login" className={styles.link}>
                    Войти
                </Link>
            </p>
        </div>
    )
}
