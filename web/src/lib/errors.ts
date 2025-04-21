import z, { ZodError } from 'zod'

export const ErrorResponseSchema = z.object({
  errorCode: z.number(),
  message: z.string(),
  readable: z.string(),
  details: z.string().optional(), // Populated on unexpected errors
})

export type ErrorResponse = z.infer<typeof ErrorResponseSchema>

export async function HandleResponseError(
  errorName: string,
  res: Response,
): Promise<ErrorResponse | null> {
  if (!res.ok) {
    try {
      const errorBody = await res.json()
      const parsedError = ErrorResponseSchema.parse(errorBody)

      console.error(
        errorName,
        `errorCode: ${parsedError.errorCode}, message: ${parsedError.message}, details: ${
          parsedError.details ?? 'No details provided'
        }`,
      )

      return parsedError
    } catch (parseError) {
      if (parseError instanceof ZodError) {
        console.error('Failed to parse error response:', parseError)

        return {
          errorCode: 500,
          message: 'Failed to parse error response',
          details: parseError.message,
          readable: 'Internal Server Error',
        }
      }
      console.error(
        'Unexpected error occurred, not instance of ZodError:',
        parseError,
      )
      return {
        errorCode: 500,
        message: 'Unexpected error occurred',
        details: String(parseError),
        readable: 'Internal Server Error',
      }
    }
  }

  return null
}
