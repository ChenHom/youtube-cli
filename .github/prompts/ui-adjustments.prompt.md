下面這份 JSON 是針對 ytcli 各命令 UI 與錯誤處理要做的「修正規格」。  
請你：

1. 讀取這份 JSON，它描述了所有需要改動的地方：  
   - transcribe 加進度條、隱藏大量 ffmpeg log  
   - chapters、overview… 加互動式選單並顯示章節數與標題  
   - 全部文字輸出支援 --lang 語言旗標  

2. 依照 JSON 裡每一個 command.ui 設定，修改對應的 Go 程式碼或 CLI 定義，並回傳更新後的程式片段。

3. 確保測試也跟著更新：例如驗證進度條顯示，驗證互動式參數正確傳遞，驗證多語系輸出。

---

```json
{
  "ui_adjustments": {
    "global": {
      "language_flag": {
        "name": "--lang",
        "type": "string",
        "default": "zh-TW",
        "description": "輸出文字說明語言，預設正體中文；可選 en、ja、…"
      }
    }
  },
  "commands": {
    "transcribe": {
      "ui": {
        "hide_verbose": true,
        "show_progress_bar": true,
        "on_error": "show_error_message_only"
      }
    },
    "chapters": {
      "ui": {
        "after_generation": {
          "show_count": true,
          "interactive_menu": {
            "type": "checkbox",
            "items_source": "chapters.json",
            "item_label_template": "Chapter {{index}}: {{title}} ({{start}}–{{end}})",
            "output_flag": "--chapters-selected"
          }
        }
      }
    },
    "overview": {
      "ui": {
        "interactive_menu": {
          "type": "radio",
          "items_source": "chapters.json",
          "item_label_template": "Chapter {{index}}: {{title}} ({{start}}–{{end}})",
          "output_flag": "--chapter"
        }
      }
    },
    "summary": {
      "ui": {
        "interactive_menu": {
          "type": "radio",
          "items_source": "chapters.json",
          "item_label_template": "Chapter {{index}}: {{title}} ({{start}}–{{end}})",
          "output_flag": "--chapter"
        }
      }
    },
    "deep-dive": {
      "ui": {
        "interactive_menu": {
          "type": "radio",
          "items_source": "chapters.json",
          "item_label_template": "Chapter {{index}}: {{title}} ({{start}}–{{end}})",
          "output_flag": "--chapter"
        }
      }
    },
    "related": {
      "ui": {
        "interactive_menu": {
          "type": "radio",
          "items_source": "chapters.json",
          "item_label_template": "Chapter {{index}}: {{title}} ({{start}}–{{end}})",
          "output_flag": "--chapter"
        }
      }
    }
  }
}
```
