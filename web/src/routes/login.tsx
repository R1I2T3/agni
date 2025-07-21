import { useTransition } from 'react'
import { FormProvider, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import type { AuthSchemaType } from '@/lib/schema'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/Input'

import { Form } from '@/components/ui/form'
import { Button } from '@/components/ui/button'
import { authSchema } from '@/lib/schema'

export const Route = createFileRoute('/login')({
  component: SignInForm,
})

function SignInForm() {
  const form = useForm<AuthSchemaType>({
    resolver: zodResolver(authSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  })
  const navigate = useNavigate()
  const [pending, startTransition] = useTransition()

  const onSubmit = (values: AuthSchemaType) => {
    startTransition(async () => {
      const res = await fetch(`/api/admin/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(values),
      })
      if (res.status !== 200) {
        toast.error('Failed to Login')
      } else {
        toast.success('Login link sent to your email')
        navigate({ to: '/' })
      }
    })
  }

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100 dark:bg-gray-900">
      <Card className="w-[90%] md:w-[50%] lg:w-[40%] mx-auto">
        <CardHeader>
          <h1 className="text-3xl font-bold m-auto">Sign In</h1>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form
              className="w-full space-y-4"
              onSubmit={form.handleSubmit(onSubmit)}
            >
              <FormProvider {...form}>
                <Input label="Username" name="username" />
                <Input label="Password" name="password" type="password" />
              </FormProvider>
              <Button
                className="w-full text-white text-[16px]"
                disabled={pending}
              >
                {pending ? 'pending...' : 'Sign In'}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  )
}

export default SignInForm
