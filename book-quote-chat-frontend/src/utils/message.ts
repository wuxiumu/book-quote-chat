// 支持多 UI 框架（如 Element Plus、Naive UI、AntD Vue、Vant），有啥用啥
// 优先 Element Plus 示例
import { ElMessage } from 'element-plus'

export const message = {
    success(msg: string) {
        ElMessage.success(msg)
    },
    info(msg: string) {
        ElMessage.info(msg)
    },
    error(msg: string) {
        ElMessage.error(msg)
    },
    warning(msg: string) {
        ElMessage.warning(msg)
    },
}