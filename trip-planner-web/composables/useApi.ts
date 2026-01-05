export const useApi = () => {
  const config = useRuntimeConfig()
  const base = config.public.apiBase as string

  const get = async <T>(path: string): Promise<T> => {
    return await $fetch<T>(`${base}${path}`)
  }

  const post = async <T>(path: string, body: any): Promise<T> => {
    return await $fetch<T>(`${base}${path}`, {
      method: 'POST',
      body,
    })
  }

  return { get, post }
}
