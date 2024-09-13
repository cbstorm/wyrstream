export default function generate(plop) {
    plop.setGenerator("module", {
      description: "Generate module",
      prompts: [
        {
          type: "input",
          name: "name",
          message: "What is the name of the module (singular)?",
        },
      ],
      actions: function (data) {
        const actions = [
          {
            type: "add",
            path: "lib/dtos/{{name}}.dto.go",
            templateFile: "__generator/__templates/dto.hbs",
          },
          {
            type: "add",
            path: "lib/repositories/{{name}}.repository.go",
            templateFile: "__generator/__templates/repository.hbs",
          },
          {
            type: "add",
            path: "lib/entities/{{name}}.entity.go",
            templateFile: "__generator/__templates/entity.hbs",
          },
          {
            type: "add",
            path: "control_service/services/{{name}}.service.go",
            templateFile: "__generator/__templates/service.hbs",
          },
          {
            type: "add",
            path: "control_service/routes/{{name}}.route.go",
            templateFile: "__generator/__templates/route.hbs",
          }
        ];
        return actions;
      },
    });
  }