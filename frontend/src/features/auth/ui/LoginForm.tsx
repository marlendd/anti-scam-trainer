import type { FormEvent } from 'react'
import {
    Link,
    useLocation,
    useNavigate,
} from 'react-router-dom'

import { useLogin } from '@/features/auth'
import { ApiError } from '@/shared/api'
import { Input, PasswordInput } from '@/shared/ui/input'
import { Logo } from '@/shared/ui/logo'

import styles from './AuthForm.module.scss'

export function LoginForm() {
    const login = useLogin()

    const navigate = useNavigate()
    const location = useLocation()

    const searchParams = new URLSearchParams(
        location.search,
    )

    const requestedPath =
        searchParams.get('returnTo')

    const returnTo =
        requestedPath?.startsWith('/') &&
        !requestedPath.startsWith('//')
            ? requestedPath
            : '/home'

    function handleSubmit(
        event: FormEvent<HTMLFormElement>,
    ) {
        event.preventDefault()

        const formData = new FormData(
            event.currentTarget,
        )

        login.mutate(
            {
                email: String(
                    formData.get('email'),
                ),
                password: String(
                    formData.get('password'),
                ),
            },
            {
                onSuccess: () => {
                    navigate(returnTo, {
                        replace: true,
                    })
                },
            },
        )
    }

    let errorMessage: string | undefined

    if (login.error instanceof ApiError) {
        if (login.error.status === 401) {
            errorMessage =
                'Неверная электронная почта или пароль'
        } else {
            errorMessage =
                'Не удалось войти. Попробуйте ещё раз'
        }
    } else if (login.isError) {
        errorMessage =
            'Произошла ошибка. Попробуйте ещё раз'
    }

    return (
        <div className={styles.content}>
            <header className={styles.header}>
                <Link
                    to="/home"
                    className={styles.brand}
                >
                    <Logo />
                </Link>

                <p
                    className={
                        styles.description
                    }
                >
                    Войдите, чтобы продолжить
                    обучение
                </p>
            </header>

            <form
                className={styles.form}
                onSubmit={handleSubmit}
            >
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
                    placeholder="Введите пароль"
                    autoComplete="current-password"
                    required
                />

                <div
                    className={
                        styles.formActions
                    }
                >
                    <Link
                        to={`/forgot-password${location.search}`}
                        className={
                            styles.secondaryLink
                        }
                    >
                        Забыли пароль?
                    </Link>
                </div>

                {errorMessage && (
                    <p className={styles.error}>
                        {errorMessage}
                    </p>
                )}

                <button
                    className={styles.submit}
                    type="submit"
                    disabled={login.isPending}
                >
                    {login.isPending
                        ? 'Входим...'
                        : 'Войти'}
                </button>
            </form>

            <p className={styles.footer}>
                Нет аккаунта?{' '}
                <Link
                    to={`/register${location.search}`}
                    className={styles.link}
                >
                    Зарегистрироваться
                </Link>
            </p>
        </div>
    )
}