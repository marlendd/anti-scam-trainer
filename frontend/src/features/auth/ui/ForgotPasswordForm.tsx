import { type FormEvent, useState } from 'react'
import { Link } from 'react-router-dom'

import { useForgot } from '@/features/auth'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Logo } from '@/shared/ui/logo'

import styles from './AuthForm.module.scss'

export function ForgotPasswordForm() {
    const [submittedEmail, setSubmittedEmail] = useState<string | null>(null)

    const forgot = useForgot()

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        const formData = new FormData(event.currentTarget)
        const email = String(formData.get('email') ?? '').trim()

        if (!email) {
            return
        }

        forgot.reset()

        try {
            await forgot.mutateAsync({
                email,
            })

            setSubmittedEmail(email)
        } catch {
            // Ошибка уже находится в forgot.error
        }
    }

    async function handleResend() {
        if (!submittedEmail) {
            return
        }

        forgot.reset()

        try {
            await forgot.mutateAsync({
                email: submittedEmail,
            })
        } catch {
            // Ошибка уже находится в forgot.error
        }
    }

    function handleRetry() {
        forgot.reset()
        setSubmittedEmail(null)
    }

    if (submittedEmail) {
        return (
            <div className={styles.content}>
                <header className={styles.header}>
                    <Link to="/home" className={styles.brand}>
                        <Logo />
                    </Link>

                    <p className={styles.description}>
                        Готово! Письмо отправили на почту <strong>{submittedEmail}</strong>
                    </p>
                </header>

                {forgot.isError && (
                    <p className={styles.error}>Не удалось отправить письмо. Попробуйте ещё раз.</p>
                )}

                <div className={styles.controls}>
                    <Button
                        type="button"
                        variant="secondary"
                        onClick={handleRetry}
                        disabled={forgot.isPending}
                    >
                        Изменить почту
                    </Button>

                    <Button type="button" onClick={handleResend} disabled={forgot.isPending}>
                        {forgot.isPending ? 'Отправляем...' : 'Отправить ещё раз'}
                    </Button>
                </div>

                <p className={styles.footer}>
                    Вспомнили пароль?{' '}
                    <Link to="/login" className={styles.link}>
                        Вернуться
                    </Link>
                </p>
            </div>
        )
    }

    return (
        <div className={styles.content}>
            <header className={styles.header}>
                <Link to="/home" className={styles.brand}>
                    <Logo />
                </Link>

                <p className={styles.description}>Введите почту, чтобы восстановить аккаунт</p>
            </header>

            <form className={styles.form} onSubmit={handleSubmit}>
                <Input
                    label="Электронная почта"
                    name="email"
                    type="email"
                    placeholder="example@mail.ru"
                    autoComplete="email"
                    inputMode="email"
                    required
                    disabled={forgot.isPending}
                />

                {forgot.isError && (
                    <p className={styles.error}>
                        Не удалось отправить письмо. Проверьте адрес и попробуйте ещё раз.
                    </p>
                )}

                <Button className={styles.submit} type="submit" disabled={forgot.isPending}>
                    {forgot.isPending ? 'Отправляем...' : 'Восстановить пароль'}
                </Button>
            </form>

            <p className={styles.footer}>
                Вспомнили пароль?{' '}
                <Link to="/login" className={styles.link}>
                    Вернуться
                </Link>
            </p>
        </div>
    )
}
