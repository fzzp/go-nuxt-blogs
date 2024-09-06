import ClipboardJS from "clipboard";

import { ElMessage } from "element-plus";

export function createBtnCopy(text: string) {
    let btnClassName = "clipboardjs_go_nuxt_blogs_copy_btn"
    let btn = document.querySelector(`.${btnClassName}`) as HTMLButtonElement;
    if (!btn) {
      btn = document.createElement("button");
      btn.style.width = "0px";
      btn.style.height = "0px";
      btn.style.display = "none";
      btn.classList.add(btnClassName)
      document.body.append(btn)
    }
  
    // 设置复制文本
    btn.setAttribute("data-clipboard-text", text)
  
    const clipboard = new ClipboardJS(`.${btnClassName}`)
    clipboard.on('error', () => {
      ElMessage.error("复制失败")
      clipboard.destroy()
    })
  
    clipboard.on('success', () => {
      ElMessage.success("复制成功")
      clipboard.destroy()
    })
  
    btn.click()
  } 