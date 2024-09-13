export default function generate(plop) {
    plop.setGenerator("dto", {
      description: "Generate Entity",
      prompts: [
        {
          type: "input",
          name: "name",
          message: "What is the name of dto (singular)?",
        },
      ],
      actions: function (data) {
        const actions = [
          {
            type: "add",
            path: "lib/dtos/{{name}}.dto.go",
            templateFile: "__generator/__templates/dto.hbs",
          },
        ];
        return actions;
      },
    });
  }