// src/shared/ui/input/PasswordInput.tsx

import { useState } from 'react'

import { Input, type InputProps } from './Input'

type PasswordInputProps = Omit<
  InputProps,
  'type' | 'endAdornment'
>

export function PasswordInput(props: PasswordInputProps) {
  const [isVisible, setIsVisible] = useState(false)

  return (
    <Input
      {...props}
      type={isVisible ? 'text' : 'password'}
      endAdornment={
        <button
          type="button"
          aria-label={
            isVisible ? 'Скрыть пароль' : 'Показать пароль'
          }
          aria-pressed={isVisible}
          onClick={() => setIsVisible((value) => !value)}
        >
          {isVisible ? 'Скрыть' : 'Показать'}
        </button>
      }
    />
  )
}