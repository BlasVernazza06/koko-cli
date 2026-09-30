import { auth } from "@[[ .ProjectName ]]/auth"

export default defineEventHandler(async (event) => {
    return auth.handler(toWebRequest(event))
})