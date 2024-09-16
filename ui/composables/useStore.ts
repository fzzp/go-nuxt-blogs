import { ref, watch } from "vue";

// localStorage的读写缓存
export function useStore<T>(key: string, initialValue: T | (() => T)) {

    if(!localStorage) { // 服务端时，没有 localStorage 会报错
        return ref<T>()
    }

    if(!import.meta.client) { // 加多一层判断，以免下面代码 localStorage 报错
        return ref<T>()
    }

    // 初始化值，如果localStorage有则取，无则使用initialValue
    const getInitValue = () => {
        const jsonValue = localStorage.getItem(key);
        if (jsonValue != null && jsonValue !== "undefined") return JSON.parse(jsonValue);
        if (typeof initialValue === 'function') {
            return (initialValue as () => T)();
        } else {
            return initialValue;
        }
    }

    // 定义一个ref值，将localStorage值缓存起来
    const sValue = ref<T>(getInitValue())

    // 监听sValue赋值时，同步到存到localStorage
    watch(
        () => sValue.value,
        (newVal) => {
            localStorage.setItem(key, JSON.stringify(newVal))
        },
        {
            deep: true
        }
    )
    
    return sValue
}
