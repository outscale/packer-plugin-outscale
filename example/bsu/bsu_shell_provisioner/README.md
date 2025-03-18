# **Outscale Packer Builder with Shell Provisioning**

This repository contains a **Packer template** for creating an **Outscale Machine Image (OMI)**. The build process provisions an Ubuntu-based VM and configures it using **shell scripts**.

---

## **📌 Overview**
This setup performs the following actions:
1. **Launches an Outscale VM** using Packer.
2. **Runs shell scripts** on the instance to install `oapi-cli` and `nginx`.
3. **Creates a new Outscale Machine Image (OMI)**.

---

## **Project Structure**
```
/bsu_shell_provisioner
├── ubuntu-oapi.pkr.hcl              # Packer template using shell provisioners
├── variables.auto.pkvars.hcl       # Variables file for Packer
│── scripts/
│   ├── install_oapi-cli.sh         # Installs Outscale oapi-cli
│── README.md                       # Documentation
```

---

## **:pushpin: Prerequisites**
Ensure the following tools are installed:

**Packer**  
   ```bash
   packer --version
   ```

---

## **:wrench: Configuration**
### **Set Up Variables (`variables.auto.pkvars.hcl`)**
Edit the `variables.auto.pkvars.hcl` file with the correct values:

```hcl
new_omi_name        = "Ubuntu-24.04-Nginx"
ssh_username        = "outscale"
vm_type             = "tinav6.c4r8p2"
region              = "eu-west-2"
osc_source_image_id = "ami-860c2495"
```

---

## **:pushpin: Running the Build**
### **Initialize Packer**
Run the following command inside the **Packer directory** to initialize required plugins:
```bash
packer init .
```

### **Validate the Configuration**
Before building, check if the configuration is valid:
```bash
packer validate .
```

### **Build the Outscale Machine Image**
```bash
packer build .
```
This process will:
Launch an **Outscale VM**.  
Run the **shell script** to install `oapi-cli`.  
Install **nginx** via `apt-get`.  
Save the **new OMI**.

---

## **Provisioning Scripts**
The **shell script (`scripts/install_oapi-cli.sh`)** ensures `oapi-cli` is installed:

Additionally, **nginx is installed** using the shell provisioner:

```hcl
provisioner "shell" {
  inline = [
    "sudo apt-get -y update",
    "sudo apt-get -y install nginx",
    "sudo systemctl enable nginx",
    "nginx -v"
  ]
}
```