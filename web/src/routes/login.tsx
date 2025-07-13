
import { createFileRoute } from '@tanstack/react-router'
import { useState } from "react"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"

export const Route = createFileRoute('/login')({
  component: Login,
})

function Login() {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        // handle login logic here
        
    }

    return (
        <div className="flex items-center justify-center min-h-screen">
            <Card className="w-full lg:w-[40dvw]">
                <CardHeader>
                    <CardTitle className='text-3xl text-center font-bold'>Login</CardTitle>
                </CardHeader>
                <form onSubmit={handleSubmit}>
                    <CardContent className="space-y-4 my-4">
                        <div>
                            <Label htmlFor="username" className='text-xl'>Username</Label>
                            <Input
                                id="username"
                                value={username}
                                onChange={e => setUsername(e.target.value)}
                                required
                                placeholder="username"
                            />
                        </div>
                        <div>
                            <Label htmlFor="password" className='text-xl'>Password</Label>
                            <Input
                                id="password"
                                type="password"
                                value={password}
                                onChange={e => setPassword(e.target.value)}
                                required
                                placeholder="••••••••"
                            />
                        </div>
                    </CardContent>
                    <CardFooter>
                        <Button type="submit" className="w-full text-xl py-2">
                            Sign In
                        </Button>
                    </CardFooter>
                </form>
            </Card>
        </div>
    )
}