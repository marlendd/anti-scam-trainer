import {
    type FormEvent,
    useState,
} from 'react'
import {Link} from 'react-router-dom'

import {Input} from '@/shared/ui/input'
import {Logo} from '@/shared/ui/logo'

import styles from './AuthForm.module.scss'
import {Button} from "@/shared/ui/button";

export function ForgotPasswordForm() {
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [submittedEmail, setSubmittedEmail] = useState<string | null>(
        null,
    )

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>,
    ) {
        event.preventDefault()

        const formData = new FormData(event.currentTarget)
        const email = String(formData.get('email') ?? '').trim()

        if (!email) {
            return
        }

        try {
            setIsSubmitting(true)

            // Здесь позже будет mutation:
            // await forgotPasswordMutation({ email }).unwrap()

            await new Promise((resolve) => {
                window.setTimeout(resolve, 700)
            })

            setSubmittedEmail(email)
        } finally {
            setIsSubmitting(false)
        }
    }

    function handleRetry() {
        setSubmittedEmail(null)
    }

    if (submittedEmail) {
        return (
            <div className={styles.content}>
                <header className={styles.header}>

                <Link to="/home" className={styles.brand}>
                    <Logo/>
                </Link>

                    <p className={styles.description}>
                        Готово! Письмо отправили на почту {' '}
                        <strong>{submittedEmail}</strong>
                    </p>
                </header>

                <div className={styles.controls}>
                    <Button
                        className={styles.submit}
                        type="button"
                        onClick={handleRetry}
                    >
                        Вернуться
                    </Button>

                    <Button
                        className={styles.submit}
                        type="button"
                        variant='secondary'
                        onClick={handleRetry}
                    >
                        Отправить ещё раз
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
                    <Logo/>
                </Link>

                <p className={styles.description}>
                    Введите почту, чтобы восстановить аккаунт
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
                    disabled={isSubmitting}
                />

                <button
                    className={styles.submit}
                    type="submit"
                    disabled={isSubmitting}
                >
                    {isSubmitting
                        ? 'Отправляем...'
                        : 'Восстановить пароль'}
                </button>
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